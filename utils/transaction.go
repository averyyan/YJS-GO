package utils

import (
	"math"
	"sort"

	"YJS-GO/structs"
	"YJS-GO/structs/content"
	"YJS-GO/types"
)

type Transaction struct {
	Doc            *YDoc
	BeforeState    map[uint64]uint64
	AfterState     map[uint64]uint64
	DeletedStructs map[*structs.Item]struct{}
	NewTypes       map[*structs.Item]struct{}
	DeleteSet      *DeleteSet
	// Changed            map[*types.AbstractType]map[string]struct{}
	Changed            map[any]map[string]struct{}
	ChangedParentTypes map[*types.AbstractType][]*YEvent
	MergeStructs       []structs.IAbstractStruct
	Origin             any
	Local              bool
	Meta               map[string]any
	SubdocsAdded       map[*YDoc]struct{}
	SubdocsRemoved     map[*YDoc]struct{}
	SubdocsLoaded      map[*YDoc]struct{}
}

func (t *Transaction) AddChangedTypeToTransaction(ty *types.AbstractType, parentSub string) {
	var item = ty.Item

	clock, ok := t.BeforeState[item.Id.Client]
	if item == nil || (ok && item.Id.Clock < clock && !item.Deleted) {
		var set map[string]struct{}
		set, ok := t.Changed[ty]
		if !ok {
			set = map[string]struct{}{}
			t.Changed[ty] = set
		}
		set[parentSub] = struct{}{}
	}
}

// RedoItem Redoes the effect of this operation.
func (t *Transaction) RedoItem(item *structs.Item, redoItems map[*structs.Item]struct{}) structs.IAbstractStruct {
	var doc = t.Doc
	var store = doc.Store
	var ownClientId = doc.ClientId
	var redone = item.Redone

	if redone != nil {
		return store.GetItemCleanStart(t, redone)
	}
	parentItem := item.Parent.(*types.AbstractType).Item

	var left structs.IAbstractStruct
	var right structs.IAbstractStruct
	if item.ParentSub == "" {
		// Is an array item. Insert at the old position.
		left = item.Left.(structs.IAbstractStruct)
		right = item
	} else {
		// Is a map item. Insert at current value.
		left = item
		for left != nil && left.(*structs.Item).Right != nil {
			left = left.(*structs.Item).Right.(*structs.Item)
			if left.ID().Client != ownClientId {
				// It is not possible to redo this item because it conflicts with a change from another client.
				return nil
			}
		}

		if left != nil && left.(*structs.Item).Right != nil {
			left = item.Parent.(*types.AbstractType).ItemMap[item.ParentSub]
		}

		right = nil
	}

	// Make sure that parent is redone.
	if parentItem != nil && parentItem.Deleted && parentItem.Redone == nil {
		// Try to undo parent if it will be undone anyway.
		if _, ok := redoItems[parentItem]; !ok || t.RedoItem(parentItem, redoItems) == nil {
			return nil
		}
	}

	if parentItem != nil && parentItem.Redone != nil {
		for parentItem.Redone != nil {
			parentItem = store.GetItemCleanStart(t, parentItem.Redone).(*structs.Item)
		}

		// Find next cloned_redo items.
		for left != nil {
			var leftTrace = left
			for leftTrace != nil && leftTrace.(*structs.Item).Parent.(*types.AbstractType).Item != parentItem {
				if leftTrace.(*structs.Item).Redone == nil {
					leftTrace = nil
				} else {
					leftTrace = store.GetItemCleanStart(t, leftTrace.(*structs.Item).Redone)
				}
			}

			if leftTrace != nil && leftTrace.(*structs.Item).Parent.(*types.AbstractType).Item == parentItem {
				left = leftTrace
				break
			}

			left = left.(*structs.Item).Left.(*structs.Item)
		}

		for right != nil {
			var rightTrace = right
			for rightTrace != nil && rightTrace.(*structs.Item).Parent.(*types.AbstractType).Item != parentItem {
				if rightTrace.(*structs.Item).Redone == nil {
					rightTrace = nil
				} else {
					rightTrace = store.GetItemCleanStart(t, rightTrace.(*structs.Item).Redone)
				}
			}

			if rightTrace != nil && rightTrace.(*structs.Item).Parent.(*types.AbstractType).Item == parentItem {
				right = rightTrace
				break
			}

			right = right.(*structs.Item).Right.(*structs.Item)
		}
	}

	var nextClock = store.GetState(ownClientId)
	var nextId = &ID{ownClientId, nextClock}
	var tmp *types.AbstractType
	if parentItem == nil {
		tmp = item.Parent.(*types.AbstractType)
	} else {
		tmp = parentItem.Content.(*content.Type).GetType()
	}
	var redoneItem = structs.NewItem(
		nextId,
		left,
		left.(*structs.Item).LastId,
		right,
		right.ID(),
		tmp,
		item.ParentSub,
		item.Content.Copy())

	item.Redone = nextId

	redoneItem.KeepItemAndParents(true)
	redoneItem.Integrate(t, 0)

	return redoneItem

}

func CleanupTransactions(transactionCleanups []*Transaction, i int) {
	if i >= len(transactionCleanups) {
		return
	}
	var transaction = transactionCleanups[i]
	var doc = transaction.Doc
	var store = doc.Store
	var ds = transaction.DeleteSet
	var mergeStructs = transaction.MergeStructs
	var actions []func()

	ds.SortAndMergeDeleteSet()
	transaction.AfterState = store.GetStateVector()
	doc.Transaction = nil

	actions = append(actions, func() {
		doc.BeforeObserverCalls(transaction)
	})
	// actions.Add(() =>
	// {
	// 	doc.InvokeOnBeforeObserverCalls(transaction)
	// })
	actions = append(actions, func() {
		for itemType, subs := range transaction.Changed {
			if itemType.(*types.AbstractType).Item == nil || !itemType.(*types.AbstractType).Item.Deleted {
				itemType.(*types.AbstractType).CallObserver(transaction, subs)
			}
		}
	})
	actions = append(actions, func() {
		// Deep observe events.
		for Type, events := range transaction.ChangedParentTypes {
			if Type.Item == nil || !Type.Item.Deleted {
				var sortedEvents []*YEvent
				for _, event := range events {
					if event.Target.Item == nil || !event.Target.Item.Deleted {
						event.CurrentTarget = Type
					}
					sortedEvents = append(sortedEvents, event)
				}
				// Sort events by path length so that top-level events are fired first.
				sort.SliceStable(sortedEvents, func(i, j int) bool {
					return len(sortedEvents[i].Path) < len(sortedEvents[j].Path)
				})
				if len(sortedEvents) <= 0 {
					continue
				}
				actions = append(actions, func() {
					Type.CallDeepEventHandlerListeners(sortedEvents, transaction)
				})
			}
		}
	})
	actions = append(actions, func() {
		doc.AfterTransaction(transaction)
	})
	for _, action := range actions {
		SafeAsyncExe(action)
	}

	// Replace deleted items with ItemDeleted / GC.
	// This is where content is actually removed from the Yjs Doc.
	if doc.GC {
		ds.TryGcDeleteSet(store, doc.GCFilter)
	}
	ds.TryMergeDeleteSet(store)

	// On all affected store.clients props, try to merge.
	for client, clock := range transaction.AfterState {
		beforeClock, ok := transaction.BeforeState[client]
		if !ok {
			beforeClock = 0
		}

		if beforeClock != clock {
			var structs = store.Clients[client]
			var firstChangePos = math.Max(float64(FindIndexSS(structs, beforeClock)), 1)
			for j := len(structs) - 1; j >= int(firstChangePos); j-- {
				TryToMergeWithLeft(structs, j)
			}
		}
	}
	// Try to merge mergeStructs.
	// TODO: It makes more sense to transform mergeStructs to a DS, sort it, and merge from right to left
	//       but at the moment DS does not handle duplicates.
	for j := 0; j < len(mergeStructs); j++ {
		var client, clock = mergeStructs[j].ID().Client, mergeStructs[j].ID().Clock
		var structs = store.Clients[client]
		var replacedStructPos = int(FindIndexSS(structs, clock))
		if replacedStructPos+1 < len(structs) {
			TryToMergeWithLeft(structs, replacedStructPos+1)
		}
		if replacedStructPos > 0 {
			TryToMergeWithLeft(structs, replacedStructPos)
		}
	}

	if !transaction.Local {
		afterClock, ok := transaction.AfterState[doc.ClientId]
		if !ok {
			afterClock = -1
		}

		beforeClock, ok := transaction.BeforeState[doc.ClientId]
		if !ok {
			beforeClock = -1
		}

		if afterClock != beforeClock {
			doc.ClientId = GenerateNewClientId()
			// Debug.WriteLine($"{nameof(Transaction)}: Changed the client-id because another client seems to be using it.");
		}
	}

	// @todo: Merge all the transactions into one and provide send the data as a single update message.
	if doc.AfterTransaction != nil {
		doc.AfterTransaction(transaction)
	}
	if doc.UpdateV2 != nil {
		doc.UpdateV2(nil, nil, transaction)
	}

	for subDoc := range transaction.SubdocsAdded {
		doc.Subdocs[subDoc] = struct{}{}
	}

	for subDoc := range transaction.SubdocsRemoved {
		delete(doc.Subdocs, subDoc)
	}
	if doc.SubdocsChanged != nil {
		doc.SubdocsChanged(transaction.SubdocsLoaded, transaction.SubdocsAdded, transaction.SubdocsRemoved)
	}
	for subDoc := range transaction.SubdocsRemoved {
		subDoc.Destroyed()
	}

	if len(transactionCleanups) <= i+1 {
		doc.TransactionCleanups = doc.TransactionCleanups[:0]
		if doc.AfterAllTransactions != nil {
			doc.AfterAllTransactions(transactionCleanups)
		}
		return
	}
	CleanupTransactions(transactionCleanups, i+1)
}

func NewTransaction(doc *YDoc, origin interface{}, local ...bool) *Transaction {
	t := &Transaction{
		Doc:                doc,
		DeleteSet:          &DeleteSet{},
		BeforeState:        doc.Store.GetStateVector(),
		AfterState:         map[uint64]uint64{},
		Changed:            map[any]map[string]struct{}{},
		ChangedParentTypes: map[*types.AbstractType][]*YEvent{},
		MergeStructs:       []structs.IAbstractStruct{},
		Meta:               map[string]any{},
		SubdocsAdded:       map[*YDoc]struct{}{},
		SubdocsRemoved:     map[*YDoc]struct{}{},
		SubdocsLoaded:      map[*YDoc]struct{}{},
		Origin:             origin,
	}
	if len(local) == 0 {
		t.Local = true
	} else {
		t.Local = false
	}
	return t
}

func SplitSnapshotAffectedStructs(transaction *Transaction, snapshot *Snapshot) {
	var (
		metaObj any
		ok      bool
	)
	if metaObj, ok = transaction.Meta["splitSnapshotAffectedStructs"]; !ok {
		metaObj = map[*Snapshot]struct{}{}
		transaction.Meta["splitSnapshotAffectedStructs"] = metaObj
	}

	var meta = metaObj.(map[*Snapshot]struct{})
	var store = transaction.Doc.Store

	// Check if we already split for this snapshot.
	if _, exist := meta[snapshot]; !exist {
		for client, clock := range snapshot.StateVector {
			if clock < store.GetState(client) {
				store.GetItemCleanStart(transaction, &ID{client, clock})
			}
		}

		snapshot.DeleteSet.IterateDeletedStructs(transaction, func(abstractStruct structs.IAbstractStruct) bool {
			return true
		})
		meta[snapshot] = struct{}{}
	}
}

// WriteUpdateMessageFromTransaction Whether the data was written.
func (t *Transaction) WriteUpdateMessageFromTransaction(encoder IUpdateEncoder) bool {
	var ok bool
	for k, v := range t.AfterState {
		if t.BeforeState[k] == v {
			ok = true
		}
	}
	if len(t.DeleteSet.Clients) == 0 && !ok {
		return false
	}
	t.DeleteSet.SortAndMergeDeleteSet()
	WriteClientsStructs(encoder, t.Doc.Store, t.BeforeState)
	t.DeleteSet.Write(encoder)
	return true
}
