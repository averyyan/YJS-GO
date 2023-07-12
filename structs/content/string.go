package content

import (
	"container/list"

	"YJS-GO/structs"
	"YJS-GO/utils"
)

var _ structs.IContentExt = (*String)(nil)

type String struct {
}

func ReadString(decoder utils.IUpdateDecoder) (String, error) {
	// TODO implement me
	panic("implement me")
}

func (s String) Copy() structs.IContent {
	// TODO implement me
	panic("implement me")
}

func (s String) Splice(offset uint64) structs.IContent {
	// TODO implement me
	panic("implement me")
}

func (s String) MergeWith(right structs.IContent) bool {
	// TODO implement me
	panic("implement me")
}

func (s String) GetContent() list.List {
	// TODO implement me
	panic("implement me")
}

func (s String) GetLength() int {
	// TODO implement me
	panic("implement me")
}

func (s String) Countable() bool {
	// TODO implement me
	panic("implement me")
}

func (s String) Write(encoder utils.IUpdateEncoder, offset int) {
	// TODO implement me
	panic("implement me")
}

func (s String) Gc(store *utils.StructStore) {
	// TODO implement me
	panic("implement me")
}

func (s String) Delete(transaction *utils.Transaction) {
	// TODO implement me
	panic("implement me")
}

func (s String) Integrate(transaction *utils.Transaction, item structs.Item) {
	// TODO implement me
	panic("implement me")
}

func (s String) GetRef() int {
	// TODO implement me
	panic("implement me")
}
