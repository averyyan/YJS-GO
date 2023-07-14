package content

import (
	"YJS-GO/structs"
	"YJS-GO/utils"
)

var _ structs.IContentExt = (*Doc)(nil)

type Doc struct {
	Doc *utils.YDoc
}

func NewDoc(doc *utils.YDoc) *Doc {
	return &Doc{doc}
}

func (d Doc) SetRef(i int) {
	// TODO implement me
	panic("implement me")
}

func ReadDoc(decoder utils.IUpdateDecoder) (Doc, error) {
	// TODO implement me
	panic("implement me")
}

func (d Doc) Copy() structs.IContent {
	// TODO implement me
	panic("implement me")
}

func (d Doc) Splice(offset uint64) structs.IContent {
	// TODO implement me
	panic("implement me")
}

func (d Doc) MergeWith(right structs.IContent) bool {
	// TODO implement me
	panic("implement me")
}

func (d Doc) GetContent() []any {
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

func (d Doc) Gc(store *utils.StructStore) {
	// TODO implement me
	panic("implement me")
}

func (d Doc) Delete(transaction *utils.Transaction) {
	// TODO implement me
	panic("implement me")
}

func (d Doc) Integrate(transaction *utils.Transaction, item *structs.Item) {
	// TODO implement me
	panic("implement me")
}

func (d Doc) GetRef() int {
	return DocRef
}
