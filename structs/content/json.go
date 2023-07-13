package content

import (
	"YJS-GO/structs"
	"YJS-GO/utils"
)

var _ structs.IContentExt = (*Json)(nil)

type Json struct {
}

func (j Json) Integrate(transaction *utils.Transaction, item *structs.Item) {
	// TODO implement me
	panic("implement me")
}

func (j Json) SetRef(i int) {
	// TODO implement me
	panic("implement me")
}

func ReadJson(decoder utils.IUpdateDecoder) (Json, error) {
	// TODO implement me
	panic("implement me")
}

func (j Json) Copy() structs.IContent {
	// TODO implement me
	panic("implement me")
}

func (j Json) Splice(offset uint64) structs.IContent {
	// TODO implement me
	panic("implement me")
}

func (j Json) MergeWith(right structs.IContent) bool {
	// TODO implement me
	panic("implement me")
}

func (j Json) GetContent() []any {
	// TODO implement me
	panic("implement me")
}

func (j Json) GetLength() int {
	// TODO implement me
	panic("implement me")
}

func (j Json) Countable() bool {
	// TODO implement me
	panic("implement me")
}

func (j Json) Write(encoder utils.IUpdateEncoder, offset int) {
	// TODO implement me
	panic("implement me")
}

func (j Json) Gc(store *utils.StructStore) {
	// TODO implement me
	panic("implement me")
}

func (j Json) Delete(transaction *utils.Transaction) {
	// TODO implement me
	panic("implement me")
}

func (j Json) GetRef() int {
	return JsonRef
}
