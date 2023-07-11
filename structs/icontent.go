package structs

import (
	"container/list"

	"YJS-GO/utils"
)

type IContent interface {
	Copy() IContent
	Splice(offset uint64) IContent
	MergeWith(right IContent) bool
	GetContent() list.List
	GetLength() int
	Countable() bool
}

type IContentExt interface {
	IContent
	Write(encoder utils.IUpdateEncoder, offset int)
	Gc(store utils.StructStore)
	Delete(transaction *utils.Transaction)
	Integrate(transaction *utils.Transaction, item Item)
	GetRef() int
	// GetType() bool
}
