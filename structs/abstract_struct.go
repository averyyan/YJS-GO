package structs

import (
	"YJS-GO/model"
	"YJS-GO/utils"
)

type AbstractStruct struct {
	Id      utils.ID
	length  int
	Deleted bool
}

type IAbstractStruct interface {
	MergeWith(AbstractStruct)

	Delete(model.Transaction)

	Integrate(transaction model.Transaction, offset int)

	GetMissing(transaction model.Transaction, store utils.StructStore)

	Write(encoder utils.IUpdateEncoder, offset int)
}
