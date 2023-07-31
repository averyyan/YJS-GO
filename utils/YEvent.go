package utils

import (
	"reflect"

	"YJS-GO/structs"
	"YJS-GO/types"
)

const (
	Add    ChangeAction = 1
	Update ChangeAction = 2
	Delete ChangeAction = 3
)

type ChangeAction int

type ChangesCollection struct {
	Added   map[*structs.Item]struct{}
	Deleted map[*structs.Item]struct{}
	Delta   []*Delta
	Keys    map[string]*ChangeKey
}

type Delta struct {
	Insert     any
	Delete     uint64
	Retain     uint64
	Attributes map[string]any
}

type ChangeKey struct {
	Action   ChangeAction
	OldValue any
}

type NewBaseType struct {
	changes       *ChangesCollection
	Changes       *ChangesCollection
	Target        *types.AbstractType
	CurrentTarget *types.AbstractType
	Transaction   *Transaction
	Path          []any
}

type YEvent struct {
	NewBaseType
	Transaction *Transaction
}

func (event *YEvent) CollectChanges() *ChangesCollection {
	if event.changes != nil {
		return event.changes
	}
	var (
		target  = event.Target
		added   = map[*structs.Item]struct{}{}
		deleted = map[*structs.Item]struct{}{}
		delta   []*Delta
		keys    = map[string]*ChangeKey{}
	)

	event.changes = &ChangesCollection{
		Added:   added,
		Deleted: deleted,
		Delta:   delta,
		Keys:    keys,
	}
	changed, ok := event.Transaction.Changed[event.Target]
	if !ok {
		changed = map[string]struct{}{}
		event.Transaction.Changed[event.Target] = changed
	}
	// if (changed.Contains(null))
	if _, ok := changed["null"]; ok { // todo 这个判断不知道怎么替换
		var lastOp *Delta
		var packOp = func() {
			if lastOp != nil {
				delta = append(delta, lastOp)
			}
		}

		for item := event.Target.Start; item != nil; item = item.Right.(*structs.Item) {
			if item.GetDeleted() {
				if event.Deletes(item) && !event.Adds(item) {
					if lastOp == nil || lastOp.Delete == 0 {
						packOp()
						lastOp = &Delta{Delete: 0}
					}

					lastOp.Delete += item.Length
					deleted[item] = struct{}{}
				} else {
					// Do nothing.
				}
			} else {
				if event.Adds(item) {
					if lastOp == nil || lastOp.Insert == nil {
						packOp()
						lastOp = &Delta{Insert: make([]any, 1)}
					}
					lastOp.Insert = lastOp.Insert.([]any)
					lastOp.Insert = append(lastOp.Insert.([]any), item.Content.GetContent().([]any)...)
					added[item] = struct{}{}
				} else {
					if lastOp == nil || lastOp.Retain == 0 {
						packOp()
						lastOp = &Delta{Retain: 0}
					}
					lastOp.Retain += item.Length
				}
			}
		}

		if lastOp != nil && lastOp.Retain == 0 {
			packOp()
		}
	}

	for key := range changed {
		if key != "" {
			var action ChangeAction
			var oldValue any
			var item = target.ItemMap[key]
			if event.Adds(item) {
				var prev = item.Left
				for prev != nil && event.Adds(prev.(*structs.Item)) {
					prev = prev.(*structs.Item).Left
				}

				if event.Deletes(item) {
					if prev != nil && event.Deletes(prev.(*structs.Item)) {
						action = Delete
						list := prev.(*structs.Item).Content.GetContent().([]any)
						oldValue = list[len(list)-1]
					} else {
						break
					}
				} else {
					if prev != nil && event.Deletes(prev.(*structs.Item)) {
						action = Update
						list := prev.(*structs.Item).Content.GetContent().([]any)
						oldValue = list[len(list)-1]
					} else {
						action = Add
						oldValue = nil
					}
				}
			} else {
				if event.Deletes(item) {
					action = Delete
					list := item.Content.GetContent().([]any)
					oldValue = list[len(list)-1]
				} else {
					break
				}
			}

			keys[key] = &ChangeKey{Action: action, OldValue: oldValue}
		}
	}
	return event.changes
}

func (event *YEvent) Deletes(item *structs.Item) bool {
	return event.Transaction.DeleteSet.IsDeleted(item.ID())
}

func (event *YEvent) Adds(item *structs.Item) bool {
	clock, ok := event.Transaction.BeforeState[item.ID().Client]
	return ok && item.ID().Clock < clock
}

func EqualAttrs(val, value any) bool {
	return reflect.DeepEqual(val, value)
}
