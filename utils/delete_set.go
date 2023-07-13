package utils

import (
	"encoding/binary"
	"math"
	"reflect"
	"sort"

	"YJS-GO/structs"
	"YJS-GO/types"
)

// DeleteSet is a temporary object that is created when needed.
// - When created in a transaction, it must only be accessed after sorting and merging.
//   - This DeleteSet is sent to other clients.
//   - We do not create a DeleteSet when we send a sync message. The DeleteSet message is created
//     directly from StructStore.
//   - We read a DeleteSet as a apart of sync/update message. In this case the DeleteSet is already
//     sorted and merged.
type DeleteSet struct {
	Clients map[uint64][]*DeleteItem
}

type DeleteItem struct {
	Clock  uint64
	Length uint64
}

func NewDeleteSet(ss *StructStore) *DeleteSet {
	ds := &DeleteSet{}
	ds.CreateDeleteSetFromStructStore(ss)
	return ds
}

func (ds *DeleteSet) CreateDeleteSetFromStructStore(ss *StructStore) {
	for k, v := range ss.Clients {
		var dsItems []*DeleteItem
		for i, str := range v {
			if str.GetDeleted() {
				clock := str.ID().Clock
				length := str.GetLength()
				for i+1 < len(v) {
					next := v[i+1]
					if next.ID().Clock == clock+length && next.GetDeleted() {
						length += next.GetLength()
						i++
					} else {
						break
					}
				}
				dsItems = append(dsItems, &DeleteItem{
					Clock:  clock,
					Length: length,
				})
			}
		}

		if len(dsItems) > 0 {
			ds.Clients[k] = dsItems
		}
	}

}

func (ds *DeleteSet) Add(client, clock, length uint64) {
	t := make([]*DeleteItem, 2)
	if item, ok := ds.Clients[client]; ok {
		t = item
	}
	t = append(t, &DeleteItem{
		Clock:  clock,
		Length: length,
	})
}

func (ds *DeleteSet) Write(encoder IDSEncoder) {
	binary.Write(encoder.RestWriter(), binary.LittleEndian, uint64(len(ds.Clients)))
	for client, dsItems := range ds.Clients {
		length := len(dsItems)
		encoder.ResetDsCurVal()
		binary.Write(encoder.RestWriter(), binary.LittleEndian, client)
		binary.Write(encoder.RestWriter(), binary.LittleEndian, length)
		for i := 0; i < length; i++ {
			item := dsItems[i]
			encoder.WriteDsClock(item.Clock)
			encoder.WriteDsLength(item.Length)
		}
	}
}

func (ds *DeleteSet) SortAndMergeDeleteSet() {
	for _, dels := range ds.Clients {
		sort.SliceStable(dels, func(i, j int) bool {
			return dels[i].Clock < dels[j].Clock
		})
		// Merge items without filtering or splicing the array.
		// i is the current pointer.
		// j refers to the current insert position for the pointed item.
		// Try to merge dels[i] into dels[j-1] or set dels[j]=dels[i].
		var i, j int
		for i, j = 1, 1; i < len(dels); i++ {
			var left = dels[j-1]
			var right = dels[i]

			if left.Clock+left.Length == right.Clock {
				dels[j-1] = &DeleteItem{
					Clock:  left.Clock,
					Length: left.Length + right.Length,
				}
				left = dels[j-1]
			} else {
				if j < i {
					dels[j] = right
				}

				j++
			}
		}
		// Trim the collection.
		if j < len(dels) {
			dels = append(dels[:j], dels[len(dels)-j:]...)
			// dels.RemoveRange(j, dels.Count-j)
		}
	}

}

func (ds *DeleteSet) TryGcDeleteSet(store *StructStore, gcFilter func(item *structs.Item) bool) {
	for client, deleteItems := range ds.Clients {
		var strs = store.Clients[client]

		for di := len(deleteItems) - 1; di >= 0; di-- {
			var deleteItem = deleteItems[di]
			var endDeleteItemClock = deleteItem.Clock + deleteItem.Length

			for si := FindIndexSS(strs, deleteItem.Clock); int(si) < len(strs); si++ {
				var str = strs[si]
				if str.ID().Clock >= endDeleteItemClock {
					break
				}
				strItem := str.(*structs.Item)
				if strItem.Deleted && !strItem.Keep && gcFilter(strItem) {
					strItem.Gc(store, false)
				}
			}
		}
	}
}

func (ds *DeleteSet) TryMergeDeleteSet(store *StructStore) {
	// Try to merge deleted / gc'd items.
	// Merge from right to left for better efficiency and so we don't miss any merge targets.
	for client, deleteItems := range ds.Clients {
		var strs = store.Clients[client]

		for di := len(deleteItems) - 1; di >= 0; di-- {
			var deleteItem = deleteItems[di]

			// Start with merging the item next to the last deleted item.
			var mostRightIndexToCheck = math.Min(float64(len(strs)-1), float64(1+FindIndexSS(strs,
				deleteItem.Clock+deleteItem.Length-1)))
			for si := int(mostRightIndexToCheck); si > 0 && strs[si].ID().Clock >= deleteItem.Clock; si-- {
				TryToMergeWithLeft(strs, si)
			}
		}
	}
}

func (ds *DeleteSet) IsDeleted(id *ID) bool {
	dis, ok := ds.Clients[id.Client]
	return ok && ds.FindIndexSS(dis, id.Clock) > 0
}

func (ds *DeleteSet) FindIndexSS(dis []*DeleteItem, clock uint64) uint64 {
	var left = 0
	var right = len(dis) - 1

	for left <= right {
		var midIndex = (left + right) / 2
		var mid = dis[midIndex]
		var midClock = mid.Clock

		if midClock <= clock {
			if clock < midClock+mid.Length {
				return uint64(midIndex)
			}

			left = midIndex + 1
		} else {
			right = midIndex - 1
		}
	}

	return 0
}

func TryToMergeWithLeft(strs []structs.IAbstractStruct, pos int) {
	var left = strs[pos-1]
	var right = strs[pos]

	if left.GetDeleted() == right.GetDeleted() && reflect.TypeOf(left).AssignableTo(reflect.TypeOf(right)) {
		if left.MergeWith(right) {
			strs = append(strs[:pos], strs[pos+1:]...)
			rightItem := right.(*structs.Item)
			if rightItem.ParentSub != "" {
				value, ok := rightItem.Parent.(types.AbstractType).ItemMap[rightItem.ParentSub]
				if ok && value == right {
					rightItem.Parent.(types.AbstractType).ItemMap[rightItem.ParentSub] = left.(*structs.Item)
				}
			}
		}
	}
}
