package content

import (
	"YJS-GO/structs"
	"YJS-GO/utils"
)

var _ structs.IContentExt = (*Binary)(nil)

type Binary struct {
	Content []byte
	Length  uint64
}

func NewBinary(content []byte) *Binary {
	return &Binary{
		content,
		1,
	}
}

func (b Binary) SetRef(i int) {
	// do nothing
}

func ReadBinary(decoder utils.IUpdateDecoder) (*Binary, error) {
	var content = decoder.ReadBuffer()
	return NewBinary(content), nil
}

func (b Binary) Copy() structs.IContent {
	return NewBinary(b.Content)
}

func (b Binary) Splice(offset uint64) structs.IContent {
	// do nothing
	return nil
}

func (b Binary) MergeWith(right structs.IContent) bool {
	return false
}

func (b Binary) GetContent() any {
	return b.Content
}

func (b Binary) GetLength() int {
	return int(b.Length)
}

func (b Binary) Countable() bool {
	return true
}

func (b Binary) Write(encoder utils.IUpdateEncoder, offset int) {
	encoder.WriteBuffer(b.Content)
}

func (b Binary) Gc(store *utils.StructStore) {
	// do nothing
}

func (b Binary) Delete(transaction *utils.Transaction) {
	// do nothing
}

func (b Binary) Integrate(transaction *utils.Transaction, item *structs.Item) {
	// do nothing
}

func (b Binary) GetRef() int {
	return BinaryRef
}
