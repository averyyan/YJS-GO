package content

import (
	"container/list"

	"YJS-GO/structs"
	"YJS-GO/utils"
)

var _ structs.IContentExt = (*Xml)(nil)

type Xml struct {
}

func ReadXml(decoder utils.IUpdateDecoder) (*Xml, error) {
	// TODO implement me
	panic("implement me")
}

func (x Xml) Copy() structs.IContent {
	// TODO implement me
	panic("implement me")
}

func (x Xml) Splice(offset uint64) structs.IContent {
	// TODO implement me
	panic("implement me")
}

func (x Xml) MergeWith(right structs.IContent) bool {
	// TODO implement me
	panic("implement me")
}

func (x Xml) GetContent() list.List {
	// TODO implement me
	panic("implement me")
}

func (x Xml) GetLength() int {
	// TODO implement me
	panic("implement me")
}

func (x Xml) Countable() bool {
	// TODO implement me
	panic("implement me")
}

func (x Xml) Write(encoder utils.IUpdateEncoder, offset int) {
	// TODO implement me
	panic("implement me")
}

func (x Xml) Gc(store utils.StructStore) {
	// TODO implement me
	panic("implement me")
}

func (x Xml) Delete(transaction *utils.Transaction) {
	// TODO implement me
	panic("implement me")
}

func (x Xml) Integrate(transaction *utils.Transaction, item structs.Item) {
	// TODO implement me
	panic("implement me")
}

func (x Xml) GetRef() int {
	// TODO implement me
	panic("implement me")
}
