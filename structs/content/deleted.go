package content

import (
	"container/list"

	"YJS-GO/structs"
	"YJS-GO/utils"
)

var _ structs.IContentExt = (*Deleted)(nil)

type Deleted struct {
}

func (d Deleted) MergeWith(right structs.IContent) bool {
	// TODO implement me
	panic("implement me")
}

func ReadDeleted(decoder utils.IUpdateDecoder) (Deleted, error) {
	// TODO implement me
	panic("implement me")
}

func (d Deleted) Copy() structs.IContent {
	// TODO implement me
	panic("implement me")
}

func (d Deleted) Splice(offset int) structs.IContent {
	// TODO implement me
	panic("implement me")
}

func (d Deleted) GetContent() list.List {
	// TODO implement me
	panic("implement me")
}

func (d Deleted) GetLength() int {
	// TODO implement me
	panic("implement me")
}

func (d Deleted) Countable() bool {
	// TODO implement me
	panic("implement me")
}

func (d Deleted) Write(encoder utils.IUpdateEncoder, offset int) {
	// TODO implement me
	panic("implement me")
}

func (d Deleted) Gc(store utils.StructStore) {
	// TODO implement me
	panic("implement me")
}

func (d Deleted) Delete(transaction utils.Transaction) {
	// TODO implement me
	panic("implement me")
}

func (d Deleted) Integrate(transaction utils.Transaction, item structs.Item) {
	// TODO implement me
	panic("implement me")
}

func (d Deleted) GetRef() int {
	// TODO implement me
	panic("implement me")
}
