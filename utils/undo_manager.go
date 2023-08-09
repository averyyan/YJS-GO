package utils

import (
	"container/list"
	"reflect"
	"time"

	"YJS-GO/structs"
	"YJS-GO/types"
)

type StackItem struct {
	BeforeState map[uint64]uint64
	AfterState  map[uint64]uint64
	Meta        map[string]any
	DeleteSet   *DeleteSet
}

func NewStackItem(ds *DeleteSet, BeforeState map[uint64]uint64,
	AfterState map[uint64]uint64) *StackItem {
	return &StackItem{
		BeforeState: BeforeState,
		AfterState:  AfterState,
		Meta:        map[string]any{},
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

// UndoManager Fires 'stack-item-added' event when a stack item was added to either the undo- or
// the redo-stack. You may store additional stack information via the metadata property
// on 'event.stackItem.meta' (it is a collection of metadata properties).
// Fires 'stack-item-popped' event when a stack item was popped from either the undo- or
// the redo-stack. You may restore the saved stack information from 'event.stackItem.Meta'.
type UndoManager struct {
	Scope           []*types.AbstractType
	DeleteFilter    func(*structs.Item) bool
	TrackedOrigins  map[any]struct{}
	UndoStack       *list.List
	RedoStack       *list.List
	Undoing         bool
	Redoing         bool
	Doc             *YDoc
	LastChange      time.Time
	CaptureTimeout  uint64
	StackItemAdded  EventHandler
	StackItemPopped EventHandler
}

func (m *UndoManager) OnAfterTransaction(transaction *Transaction) {
	// Only track certain transactions.
	var exist bool
	for _, abstractType := range m.Scope {
		_, ok := transaction.ChangedParentTypes[abstractType]
		if ok {
			exist = true
			break
		}
	}
	var assignable bool
	for to := range m.TrackedOrigins {
		if reflect.TypeOf(to).AssignableTo(reflect.TypeOf(transaction.Origin)) {
			assignable = true
			break
		}
	}

	if _, tmpOk := m.TrackedOrigins[transaction.Origin]; !exist || tmpOk && (transaction.Origin == nil || !assignable) {
		return
	}

	var undoing, redoing = m.Undoing, m.Redoing
	var stack = m.UndoStack
	if undoing {
		stack = m.RedoStack
	}

	if undoing {
		// Next undo should not be appended to last stack item.
		m.stopCapturing()
	} else if !redoing {
		// Neither undoing nor redoing: delete redoStack.
		m.RedoStack = list.New()
	}

	var beforeState, afterState = transaction.BeforeState, transaction.AfterState

	var now = time.Now()
	if time.Now().Sub(m.LastChange).Milliseconds() < int64(m.CaptureTimeout) && stack.Len() > 0 && !undoing && !redoing {
		// Append change to last stack op.
		var lastOp = stack.Back().Value.(*StackItem)
		lastOp.DeleteSet = NewDeleteSetWithArray([]*DeleteSet{lastOp.DeleteSet, transaction.DeleteSet})
		lastOp.AfterState = afterState
	} else {
		// Create a new stack op.
		var item = NewStackItem(transaction.DeleteSet, beforeState, afterState)
		stack.PushBack(item)
	}
	if !undoing && !redoing {
		m.LastChange = now
	}
	// Make sure that deleted structs are not GC'd.
	transaction.DeleteSet.IterateDeletedStructs(transaction, func(i structs.IAbstractStruct) bool {
		item, ok := i.(*structs.Item)
		var exist bool
		for _, abstractType := range m.Scope {
			if IsParentOf(abstractType, item) {
				exist = true
				break
			}
		}
		if ok && exist {
			item.KeepItemAndParents(true)
		}
		return true
	})

	if m.StackItemAdded != nil {
		opType := undo
		if undoing {
			opType = redo
		}
		m.StackItemAdded(NewStackEventArgs(stack.Back().Value.(*StackItem), opType, transaction.ChangedParentTypes, transaction.Origin))
	}
}

func NewUndoManager(typeScopes []*types.AbstractType, captureTimeout uint64,
	deleteFilter func(*structs.Item) bool,
	trackedOrigins map[any]struct{}) *UndoManager {
	t := &UndoManager{
		Scope:          typeScopes,
		DeleteFilter:   func(item *structs.Item) bool { return true },
		TrackedOrigins: map[any]struct{}{},
		UndoStack:      list.New(),
		RedoStack:      list.New(),
		Undoing:        false,
		Redoing:        false,
		Doc:            typeScopes[0].Doc,
		LastChange:     time.Now().Add(-1 * 100 * 12 * 31 * 24 * time.Hour),
		CaptureTimeout: captureTimeout,
	}
	if deleteFilter != nil {
		t.DeleteFilter = deleteFilter
	}
	if trackedOrigins != nil {
		t.TrackedOrigins = trackedOrigins
	}
	t.TrackedOrigins[t] = struct{}{}
	t.Doc.AfterTransaction = t.OnAfterTransaction
	return t
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
							break
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
func (m *UndoManager) stopCapturing() {
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

func (m *UndoManager) Redo() *StackItem {
	m.Redoing = true

	if res, err := m.PopStackItem(m.RedoStack, redo); err != nil {
		m.Redoing = false
	} else {
		return res
	}
	return nil
}

func (m *UndoManager) PopStackItem(stack *list.List, operationType OperationType) (*StackItem, error) {
	var (
		result *StackItem
		tr     *Transaction
	)
	// Keep a reference to the transaction, so we can fire the event with the 'changedParentTypes'.
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

				var clockLen = endClock - startClock
				var strs = m.Doc.Store.Clients[client]
				if startClock == endClock {
					continue
				}
				// Make sure strs don't overlap with the range of created operations [stackItem.start, stackItem.start + stackItem.end).
				// This must be executed before deleted strs are iterated.
				m.Doc.Store.GetItemCleanStart(transaction, &ID{client, startClock})
				if endClock < m.Doc.Store.GetState(client) {
					m.Doc.Store.GetItemCleanStart(transaction, &ID{client, startClock})
				}
				m.Doc.Store.IterateStructs(transaction, strs, startClock, clockLen, func(str structs.IAbstractStruct) bool {
					if it, ok := str.(*structs.Item); ok {
						if it.Redone != nil {
							var item, diff = m.Doc.Store.FollowRedone(str.ID())
							if diff > 0 {
								item = m.Doc.Store.GetItemCleanStart(transaction, &ID{item.ID().Client,
									item.ID().Clock + diff})
							}
							if item.GetLength() > clockLen {
								m.Doc.Store.GetItemCleanStart(transaction, &ID{item.ID().Client, endClock})
							}
							str, it = item.(*structs.Item), item.(*structs.Item)
						}
						var exist bool
						for _, abstractType := range m.Scope {
							if IsParentOf(abstractType, it) {
								exist = true
								break
							}
						}
						if !it.Deleted && exist {
							itemsToDelete = append(itemsToDelete, it)
						}
					}
					return true
				})

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
						break
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
