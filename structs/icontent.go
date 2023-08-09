package structs

import (
	"YJS-GO/utils"
)

type IContent interface {
	Copy() IContentExt
	Splice(offset uint64) IContentExt
	MergeWith(right IContentExt) bool
	GetContent() any
	GetLength() int
	Countable() bool
}

type IContentExt interface {
	IContent
	Write(encoder utils.IUpdateEncoder, offset int)
	Gc(store *utils.StructStore)
	Delete(transaction *utils.Transaction)
	Integrate(transaction *utils.Transaction, item *Item)
	GetRef() int
	SetRef(int)
	// GetType() bool
}
