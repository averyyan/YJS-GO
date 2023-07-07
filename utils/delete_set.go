package utils

import "encoding/binary"

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

func (ds DeleteSet) CreateDeleteSetFromStructStore(ss *StructStore) {
	for k, v := range ss.Clients {
		var dsItems []*DeleteItem

		for i, str := range v {
			if str.GetDeleted() {
				clock := str.ID().Clock
				length := str.GetLength()
				for i+1 < len(v) {
					next := v[i+1]
					if next.ID().Clock == clock+uint64(length) && next.GetDeleted() {
						length += next.GetLength()
						i++
					} else {
						break
					}
				}
				dsItems = append(dsItems, &DeleteItem{
					Clock:  clock,
					Length: uint64(length),
				})
			}
		}

		if len(dsItems) > 0 {
			ds.Clients[k] = dsItems
		}
	}

}

func (ds DeleteSet) Add(client, clock, length uint64) {
	t := make([]*DeleteItem, 2)
	if item, ok := ds.Clients[client]; ok {
		t = item
	}
	t = append(t, &DeleteItem{
		Clock:  clock,
		Length: length,
	})
}

func (ds DeleteSet) Write(encoder IDSEncoder) {
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
