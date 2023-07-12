package content

import (
	"container/list"

	"YJS-GO/structs"
	"YJS-GO/utils"
)

var _ structs.IContentExt = (*Any)(nil)

type Any struct {
}

func ReadAny(decoder utils.IUpdateDecoder) (Any, error) {
	// TODO implement me
	panic("implement me")
}

func (a Any) Copy() structs.IContent {
	// TODO implement me
	panic("implement me")
}

func (a Any) Splice(offset uint64) structs.IContent {
	// TODO implement me
	panic("implement me")
}

func (a Any) MergeWith(right structs.IContent) bool {
	// TODO implement me
	panic("implement me")
}

func (a Any) GetContent() list.List {
	// TODO implement me
	panic("implement me")
}

func (a Any) GetLength() int {
	// TODO implement me
	panic("implement me")
}

func (a Any) Countable() bool {
	// TODO implement me
	panic("implement me")
}

func (a Any) Write(encoder utils.IUpdateEncoder, offset int) {
	// TODO implement me
	panic("implement me")
}

func (a Any) Gc(store *utils.StructStore) {
	// TODO implement me
	panic("implement me")
}

func (a Any) Delete(transaction *utils.Transaction) {
	// TODO implement me
	panic("implement me")
}

func (a Any) Integrate(transaction *utils.Transaction, item structs.Item) {
	// TODO implement me
	panic("implement me")
}

func (a Any) GetRef() int {
	// TODO implement me
	panic("implement me")
}
