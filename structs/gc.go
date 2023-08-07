package structs

import "YJS-GO/utils"

var _ IAbstractStruct = (*GC)(nil)

const StructGCRefNumber = 0

type GC struct {
	Id      *utils.ID
	Length  uint64
	Deleted bool
}

func NewGC(id *utils.ID, length uint64) *GC {
	return &GC{
		Id:     id,
		Length: length,
	}
}

func (G *GC) GetLength() uint64 {
	return G.Length
}

func (G *GC) GetDeleted() bool {
	return G.Deleted
}

func (G *GC) ID() *utils.ID {
	return G.Id
}

func (G *GC) MergeWith(right any) bool {
	if _, ok := right.(*GC); !ok {
		return false
	}
	// Debug.Assert(right is GC)
	G.Length += right.(*GC).GetLength()
	return true
}

func (G *GC) Delete(transaction *utils.Transaction) {
	// Do nothing.
}

func (G *GC) Integrate(transaction *utils.Transaction, offset int) {
	if offset > 0 {
		G.Id = &utils.ID{Client: G.Id.Client, Clock: G.Id.Clock + uint64(offset)}
		G.Length -= uint64(offset)
	}

	transaction.Doc.Store.AddStruct(G)
}

func (G *GC) GetMissing(transaction *utils.Transaction, store *utils.StructStore) (uint64, error) {
	return 0, nil
}

func (G *GC) Write(encoder utils.IUpdateEncoder, offset int) {
	encoder.WriteInfo(StructGCRefNumber)
	encoder.WriteLength(int(G.Length) - offset)
}
