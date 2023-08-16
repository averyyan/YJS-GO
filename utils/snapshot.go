package utils

import "encoding/binary"

type Snapshot struct {
	DeleteSet   DeleteSet
	StateVector map[uint64]uint64
}

func NewSnapshot(ds DeleteSet, stateMap map[uint64]uint64) *Snapshot {
	return &Snapshot{
		DeleteSet:   ds,
		StateVector: stateMap,
	}
}

func (s *Snapshot) Equal(other *Snapshot) bool {
	if other == nil {
		return false
	}
	var ds1 = s.DeleteSet.Clients
	var ds2 = other.DeleteSet.Clients
	var sv1 = s.StateVector
	var sv2 = other.StateVector
	// todo 取指针对比可能有问题，可能需要修复
	if &ds1 != &ds2 || &sv1 != &sv2 {
		return false
	}
	for k, v := range sv1 {
		v2, ok := sv2[k]
		if !ok {
			return false
		}
		if v2 != v {
			return false
		}
	}
	for k, v := range ds1 {
		i, ok := ds2[k]
		if !ok {
			return false
		}
		if len(i) != len(v) {
			return false
		}
		for i2, item := range i {
			var dsItem2 = v[i2]
			if item.Clock != dsItem2.Clock || item.Length != dsItem2.Length {
				return false
			}
		}
	}
	return true
}

func (s *Snapshot) EncodeSnapshotV2() []byte {
	var encoder = NewDsDecoderV2WithEmpty()

	s.DeleteSet.Write(encoder)
	WriteStateVector(encoder, s.StateVector)
	return encoder.ToArray()

}

func (s *Snapshot) RestoreDocument(originDoc *YDoc, opts *YDocOptions) *YDoc {
	if originDoc.GC {
		// We should try to restore a GC-ed document, because some of the restored items might have their content deleted.
		// throw new Exception("originDoc must not be garbage collected");
		return nil
	}
	encoder := NewUpdateEncoderV2()
	originDoc.Transact(func(tr *Transaction) {
		count := 0
		for _, v := range s.StateVector {
			if v > 0 {
				count++
			}
		}
		size := count
		binary.Write(encoder.RestWriter(), binary.LittleEndian, size)

		// Splitting the structs before writing them to the encoder.
		for client, clock := range s.StateVector {
			if clock == 0 {
				continue
			}
			if clock < originDoc.Store.GetState(client) {
				tr.Doc.Store.GetItemCleanStart(tr, &ID{
					Client: client,
					Clock:  clock,
				})
			}

			var structs = originDoc.Store.Clients[client]
			var lastStructIndex = FindIndexSS(structs, clock-1)

			// Write # encoded structs.
			binary.Write(encoder.RestWriter(), binary.LittleEndian, uint64(lastStructIndex+1))
			encoder.WriteClient(client)

			// First clock written is 0.
			binary.Write(encoder.RestWriter(), binary.LittleEndian, 0)

			for i := 0; i < int(lastStructIndex); i++ {
				structs[i].Write(encoder, 0)
			}
		}

		s.DeleteSet.Write(encoder)
	}, nil)

	var newDoc *YDoc
	if opts != nil {
		newDoc = NewDoc(opts)
	} else {
		newDoc = NewDoc(originDoc.CloneOptionsWithNewGuid())
	}
	newDoc.ApplyUpdateV2(encoder.ToArray(), "snapshot")
	return newDoc
}
