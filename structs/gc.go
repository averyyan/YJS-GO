package structs

import "YJS-GO/utils"

var _ IAbstractStruct = (*GC)(nil)

type GC struct {
	Id      utils.ID
	Length  uint64
	Deleted bool
}

func (G GC) GetLength() uint {
	return uint(G.Length)
}

func (G GC) GetDeleted() bool {
	return G.Deleted
}

func (G GC) ID() utils.ID {
	return G.Id
}

func (G GC) MergeWith(a any) bool {
	// TODO implement me
	panic("implement me")
}

func (G GC) Delete(transaction utils.Transaction) {
	// TODO implement me
	panic("implement me")
}

func (G GC) Integrate(transaction utils.Transaction, offset int) {
	// TODO implement me
	panic("implement me")
}

func (G GC) GetMissing(transaction utils.Transaction, store utils.StructStore) {
	// TODO implement me
	panic("implement me")
}

func (G GC) Write(encoder utils.IUpdateEncoder, offset int) {
	// TODO implement me
	panic("implement me")
}
