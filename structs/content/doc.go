package content

import (
	"container/list"

	"YJS-GO/structs"
	"YJS-GO/utils"
)

var _ structs.IContentExt = (*Doc)(nil)

type Doc struct {
}

func ReadDoc(decoder utils.IUpdateDecoder) (Doc, error) {
	// TODO implement me
	panic("implement me")
}

func (d Doc) Copy() structs.IContent {
	// TODO implement me
	panic("implement me")
}

func (d Doc) Splice(offset int) structs.IContent {
	// TODO implement me
	panic("implement me")
}

func (d Doc) MergeWith(right structs.IContent) bool {
	// TODO implement me
	panic("implement me")
}

func (d Doc) GetContent() list.List {
	// TODO implement me
	panic("implement me")
}

func (d Doc) GetLength() int {
	// TODO implement me
	panic("implement me")
}

func (d Doc) Countable() bool {
	// TODO implement me
	panic("implement me")
}

func (d Doc) Write(encoder utils.IUpdateEncoder, offset int) {
	// TODO implement me
	panic("implement me")
}

func (d Doc) Gc(store utils.StructStore) {
	// TODO implement me
	panic("implement me")
}

func (d Doc) Delete(transaction utils.Transaction) {
	// TODO implement me
	panic("implement me")
}

func (d Doc) Integrate(transaction utils.Transaction, item structs.Item) {
	// TODO implement me
	panic("implement me")
}

func (d Doc) GetRef() int {
	// TODO implement me
	panic("implement me")
}
