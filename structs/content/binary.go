package content

import (
	"YJS-GO/structs"
	"YJS-GO/utils"
)

var _ structs.IContentExt = (*Binary)(nil)

type Binary struct {
}

func (b Binary) SetRef(i int) {
	// TODO implement me
	panic("implement me")
}

func ReadBinary(decoder utils.IUpdateDecoder) (Binary, error) {
	// TODO implement me
	panic("implement me")
}

func (b Binary) Copy() structs.IContent {
	// TODO implement me
	panic("implement me")
}

func (b Binary) Splice(offset uint64) structs.IContent {
	// TODO implement me
	panic("implement me")
}

func (b Binary) MergeWith(right structs.IContent) bool {
	// TODO implement me
	panic("implement me")
}

func (b Binary) GetContent() []any {
	// TODO implement me
	panic("implement me")
}

func (b Binary) GetLength() int {
	// TODO implement me
	panic("implement me")
}

func (b Binary) Countable() bool {
	// TODO implement me
	panic("implement me")
}

func (b Binary) Write(encoder utils.IUpdateEncoder, offset int) {
	// TODO implement me
	panic("implement me")
}

func (b Binary) Gc(store *utils.StructStore) {
	// TODO implement me
	panic("implement me")
}

func (b Binary) Delete(transaction *utils.Transaction) {
	// TODO implement me
	panic("implement me")
}

func (b Binary) Integrate(transaction *utils.Transaction, item *structs.Item) {
	// TODO implement me
	panic("implement me")
}

func (b Binary) GetRef() int {
	return BinaryRef
}
