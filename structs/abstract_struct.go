package structs

import (
	"YJS-GO/utils"
)

type AbstractStruct struct {
	Id      utils.ID
	Length  int
	Deleted bool
}

type IAbstractStruct interface {
	MergeWith(any) bool
	Delete(utils.Transaction)
	Integrate(transaction utils.Transaction, offset int)
	GetMissing(transaction utils.Transaction, store utils.StructStore)
	Write(encoder utils.IUpdateEncoder, offset int)
}
