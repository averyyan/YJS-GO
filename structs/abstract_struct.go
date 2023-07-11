package structs

import (
	"YJS-GO/utils"
)

type AbstractStruct struct {
	Id      *utils.ID
	Length  uint64
	Deleted bool
}

type IAbstractStruct interface {
	MergeWith(any) bool
	Delete(*utils.Transaction)
	Integrate(transaction *utils.Transaction, offset int)
	GetMissing(transaction *utils.Transaction, store *utils.StructStore) (uint64, error)
	Write(encoder utils.IUpdateEncoder, offset int)
	ID() *utils.ID
	GetLength() uint64
	GetDeleted() bool
}
