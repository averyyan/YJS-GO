package content

import (
	"YJS-GO/structs"
	"YJS-GO/types"
	"YJS-GO/utils"
)

var _ structs.IContentExt = (*Format)(nil)

type Format struct {
	Key   string
	Value any
}

func NewFormat(key string, value any) *Format {
	return &Format{key, value}
}

func (f Format) SetRef(i int) {
	// Do nothing.

}

func ReadFormat(decoder utils.IUpdateDecoder) (*Format, error) {
	var key = decoder.ReadKey()
	var value = decoder.ReadJson()
	return NewFormat(key, value), nil
}

func (f Format) Copy() structs.IContentExt {
	return NewFormat(f.Key, f.Value)
}

func (f Format) Splice(offset uint64) structs.IContentExt {
	// Do nothing.
	return nil
}

func (f Format) MergeWith(right structs.IContentExt) bool {
	return false
}

func (f Format) GetContent() any {
	// Do nothing.
	return nil
}

func (f Format) GetLength() int {
	return 1
}

func (f Format) Countable() bool {
	return false
}

func (f Format) Write(encoder utils.IUpdateEncoder, offset int) {
	encoder.WriteKey(f.Key)
	encoder.WriteJson(f.Value)
}

func (f Format) Gc(store *utils.StructStore) {
	// Do nothing.

}

func (f Format) Delete(transaction *utils.Transaction) {
	// Do nothing.

}

func (f Format) Integrate(transaction *utils.Transaction, item *structs.Item) {
	// Search markers are currently unsupported for rich text documents.
	item.Parent.(*types.YArrayBase).ClearSearchMarkers()
}

func (f Format) GetRef() int {
	return FormatRef
}
