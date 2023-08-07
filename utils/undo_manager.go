package utils

import (
	"container/list"
	"time"

	"YJS-GO/structs"
	"YJS-GO/types"
)

type StackItem struct {
	BeforeState map[uint64]uint64
	AfterState  map[uint64]uint64
	Meta        map[string]interface{}
	DeleteSet   DeleteSet
}

func NewStackItem(BeforeState map[uint64]uint64,
	AfterState map[uint64]uint64,
	Meta map[string]interface{},
	ds DeleteSet) *StackItem {
	return &StackItem{
		BeforeState: BeforeState,
		AfterState:  AfterState,
		Meta:        Meta,
		DeleteSet:   ds,
	}
}

// OperationType 操作类型
type OperationType string

const (
	undo = OperationType("undo")
	redo = OperationType("redo")
)

type EventHandler func(v any)

type StackEventArgs struct {
	StackItem          *StackItem
	Type               OperationType
	ChangedParentTypes map[*types.AbstractType][]*YEvent
	Origin             any
}

type UndoManager struct {
	Scope          []*types.AbstractType
	DeleteFilter   func(*structs.Item) bool
	TrackedOrigins map[any]struct {
	}
	UndoStack       *list.List
	RedoStack       *list.List
	Undoing         bool
	Redoing         bool
	Doc             YDoc
	LastChange      time.Time
	CaptureTimeout  uint64
	StackItemAdded  EventHandler
	StackItemPopped EventHandler
}

func NewUndoManager() *UndoManager {
	return &UndoManager{
		Scope:          nil,
		DeleteFilter:   nil,
		TrackedOrigins: nil,
		UndoStack:      list.New(),
		RedoStack:      list.New(),
		Undoing:        false,
		Redoing:        false,
		Doc:            YDoc{},
		LastChange:     time.Time{},
		CaptureTimeout: 0,
	}
}

func (m *UndoManager) GetCount() uint64 {
	return uint64(m.UndoStack.Len())
}

func (m *UndoManager) Clear() {
	m.Doc.Transact(func(tr *Transaction) {
		clearItem := func(stackItem any) {
			stackItem.(*StackItem).DeleteSet.IterateDeletedStructs(tr,
				func(i structs.IAbstractStruct) bool {
					item, ok := i.(*structs.Item)
					var exist bool
					for _, abstractType := range m.Scope {
						if IsParentOf(abstractType, item) {
							exist = true
						}
					}
					if ok && exist {
						item.KeepItemAndParents(false)
					}
					return true
				})
		}
		for item := m.UndoStack.Front(); item != nil; {
			clearItem(item)
		}

		for item := m.RedoStack.Front(); item != nil; {
			clearItem(item)
		}
	}, nil)

	m.UndoStack = list.New()
	m.RedoStack = list.New()
}

// StopCapturing UndoManager merges Undo-StackItem if they are created within time-gap
// smaller than 'captureTimeout'. Call this method so that the next StackItem
// won't be merged.
func (m *UndoManager) StopCapturing() {
	//            _lastChange = DateTime.MinValue;
	m.LastChange = time.Now().Add(-1 * 100 * 12 * 31 * 24 * time.Hour)
}

func (m *UndoManager) Undo() *StackItem {
	m.Undoing = true

	if res, err := m.PopStackItem(m.UndoStack, undo); err != nil {
		m.Undoing = false
	} else {
		return res
	}
	return nil
}

func (m *UndoManager) PopStackItem(stack *list.List, operationType OperationType) (*StackItem, error) {
	var result *StackItem

	// Keep a reference to the transaction so we can fire the event with the 'changedParentTypes'.
	var tr *Transaction

	m.Doc.Transact(func(transaction *Transaction) {
		tr = transaction

		for stack.Len() > 0 && result == nil {
			var stackItem = stack.Back().Value.(*StackItem)
			var itemsToRedo = map[*structs.Item]struct{}{}
			var itemsToDelete []*structs.Item
			var performedChange = false
			for client, endClock := range stackItem.AfterState {
				var (
					startClock uint64
					ok         bool
				)
				if startClock, ok = stackItem.BeforeState[client]; !ok {
					startClock = 0
				}

				var len = endClock - startClock
				var strs = m.Doc.Store.Clients[client]

				if startClock != endClock {
					// Make sure strs don't overlap with the range of created operations [stackItem.start, stackItem.start + stackItem.end).
					// This must be executed before deleted strs are iterated.
					m.Doc.Store.GetItemCleanStart(transaction, &ID{client, startClock})

					if endClock < m.Doc.Store.GetState(client) {
						m.Doc.Store.GetItemCleanStart(transaction, &ID{client, startClock})
					}

					m.Doc.Store.IterateStructs(transaction, strs, startClock, len, func(str structs.IAbstractStruct) bool {
						if it, ok := str.(*structs.Item); ok {
							if it.Redone != nil {
								var item, diff = m.Doc.Store.FollowRedone(str.ID())

								if diff > 0 {
									item = m.Doc.Store.GetItemCleanStart(transaction, &ID{item.ID().Client,
										item.ID().Clock + diff})
								}

								if item.GetLength() > len {
									m.Doc.Store.GetItemCleanStart(transaction, &ID{item.ID().Client, endClock})
								}

								str, it = item.(*structs.Item), item.(*structs.Item)
							}
							var exist bool
							for _, abstractType := range m.Scope {
								if IsParentOf(abstractType, it) {
									exist = true
								}
							}
							if !it.Deleted && exist {
								itemsToDelete = append(itemsToDelete, it)
							}
						}
						return true
					})
				}
			}

			stackItem.DeleteSet.IterateDeletedStructs(transaction, func(str structs.IAbstractStruct) bool {
				var id = str.ID()
				var clock = id.Clock
				var client = id.Client

				var (
					startClock uint64
					ok         bool
				)
				if startClock, ok = stackItem.BeforeState[client]; !ok {
					startClock = 0
				}
				var endClock uint64
				if endClock, ok = stackItem.AfterState[client]; !ok {
					endClock = 0
				}
				item, ok := str.(*structs.Item)
				var exist bool
				for _, abstractType := range m.Scope {
					if IsParentOf(abstractType, item) {
						exist = true
					}
				}
				// Never redo structs in [stackItem.start, stackItem.start + stackItem.end), because they were created and deleted in the same capture interval.
				if ok && exist && !(clock >= startClock && clock < endClock) {
					itemsToRedo[item] = struct{}{}
				}
				return true
			})
			for str, _ := range itemsToRedo {
				tmpBool := transaction.RedoItem(str, itemsToRedo) == nil
				performedChange = performedChange || tmpBool
			}

			// We want to delete in reverse order so that children are deleted before
			// parents, so we have more information available when items are filtered.
			for i := len(itemsToDelete) - 1; i >= 0; i-- {
				var item = itemsToDelete[i]
				if m.DeleteFilter(item) {
					item.Delete(transaction)
					performedChange = true
				}
			}

			result = stackItem
		}
		for Type, subProps := range transaction.Changed {
			// Destroy search marker if necessary.
			arr, ok := Type.(*types.YArrayBase)
			if _, exist := subProps[""]; exist && ok {
				arr.ClearSearchMarkers()
			}
		}
	}, m)

	if result != nil {
		m.StackItemPopped(NewStackEventArgs(result, operationType, tr.ChangedParentTypes, tr.Origin))
	}

	return result, nil
}

func NewStackEventArgs(result *StackItem, eventType OperationType, parentTypes map[*types.AbstractType][]*YEvent, origin any) interface{} {
	return StackEventArgs{
		StackItem:          result,
		Type:               eventType,
		ChangedParentTypes: parentTypes,
		Origin:             origin,
	}
}

func IsParentOf(parent *types.AbstractType, child *structs.Item) bool {
	for child != nil {
		if child.Parent == parent {
			return true
		}
		child = child.Parent.(*types.AbstractType).Item
	}

	return false
}
