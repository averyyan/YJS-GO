package content

import (
	"container/list"

	"YJS-GO/structs"
	"YJS-GO/utils"
)

var _ structs.IContentExt = (*Type)(nil)

type Type struct {
}

func ReadType(decoder utils.IUpdateDecoder) (Type, error) {
	// TODO implement me
	panic("implement me")
}

func (t Type) Copy() structs.IContent {
	// TODO implement me
	panic("implement me")
}

func (t Type) Splice(offset uint64) structs.IContent {
	// TODO implement me
	panic("implement me")
}

func (t Type) MergeWith(right structs.IContent) bool {
	// TODO implement me
	panic("implement me")
}

func (t Type) GetContent() list.List {
	// TODO implement me
	panic("implement me")
}

func (t Type) GetLength() int {
	// TODO implement me
	panic("implement me")
}

func (t Type) Countable() bool {
	// TODO implement me
	panic("implement me")
}

func (t Type) Write(encoder utils.IUpdateEncoder, offset int) {
	// TODO implement me
	panic("implement me")
}

func (t Type) Gc(store utils.StructStore) {
	// TODO implement me
	panic("implement me")
}

func (t Type) Delete(transaction *utils.Transaction) {
	// TODO implement me
	panic("implement me")
}

func (t Type) Integrate(transaction *utils.Transaction, item structs.Item) {
	// TODO implement me
	panic("implement me")
}

func (t Type) GetRef() int {
	// TODO implement me
	panic("implement me")
}
