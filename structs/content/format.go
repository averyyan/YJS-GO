package content

import (
	"container/list"

	"YJS-GO/structs"
	"YJS-GO/utils"
)

var _ structs.IContentExt = (*Format)(nil)

type Format struct {
}

func ReadFormat(decoder utils.IUpdateDecoder) (Format, error) {
	// TODO implement me
	panic("implement me")
}

func (f Format) Copy() structs.IContent {
	// TODO implement me
	panic("implement me")
}

func (f Format) Splice(offset uint64) structs.IContent {
	// TODO implement me
	panic("implement me")
}

func (f Format) MergeWith(right structs.IContent) bool {
	// TODO implement me
	panic("implement me")
}

func (f Format) GetContent() list.List {
	// TODO implement me
	panic("implement me")
}

func (f Format) GetLength() int {
	// TODO implement me
	panic("implement me")
}

func (f Format) Countable() bool {
	// TODO implement me
	panic("implement me")
}

func (f Format) Write(encoder utils.IUpdateEncoder, offset int) {
	// TODO implement me
	panic("implement me")
}

func (f Format) Gc(store utils.StructStore) {
	// TODO implement me
	panic("implement me")
}

func (f Format) Delete(transaction *utils.Transaction) {
	// TODO implement me
	panic("implement me")
}

func (f Format) Integrate(transaction *utils.Transaction, item structs.Item) {
	// TODO implement me
	panic("implement me")
}

func (f Format) GetRef() int {
	// TODO implement me
	panic("implement me")
}
