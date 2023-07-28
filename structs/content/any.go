package content

import (
	"YJS-GO/structs"
	"YJS-GO/utils"
)

var _ structs.IContentExt = (*Any)(nil)

type Any struct {
	Ref     int
	Content []any
	Cable   bool
	Length  int
}

func (a Any) SetRef(ref int) {
	a.Ref = ref
}

func NewAny(content any) *Any {
	return &Any{
		Ref:     AnyRef,
		Content: content.([]any),
		Cable:   true,
		Length:  len(content.([]any)),
	}
}

func ReadAny(decoder utils.IUpdateDecoder) (**Any, error) {
	var length = decoder.ReadLength()
	var cs = make([]any, length)

	for i := 0; i < int(length); i++ {
		var c = decoder.ReadAny()
		cs = append(cs, c)
	}

	return NewAny(cs), nil
}

func (a Any) Copy() structs.IContent {
	return NewAny(a.Content)
}

func (a Any) Splice(offset uint64) structs.IContent {
	var right = NewAny(a.Content[offset : len(a.Content)-int(offset)])
	a.Content = append(a.Content[:offset], a.Content[len(a.Content)-int(offset):])
	return right
}

func (a Any) MergeWith(right structs.IContent) bool {
	a.Content = append(a.Content, (right.(*Any)).Content)
	return true
}

func (a Any) GetContent() any {
	return a.Content
}

func (a Any) GetLength() int {
	return a.Length
}

func (a Any) Countable() bool {
	return a.Cable
}

func (a Any) Write(encoder utils.IUpdateEncoder, offset int) {
	length := len(a.Content)
	encoder.WriteLength(length - offset)

	for i := offset; i < length; i++ {
		var c = a.Content[i]
		encoder.WriteAny(c)
	}
}

func (a Any) Gc(store *utils.StructStore) {
	// Do nothing.
}

func (a Any) Delete(transaction *utils.Transaction) {
	// Do nothing.
}

func (a Any) Integrate(transaction *utils.Transaction, item *structs.Item) {
	// Do nothing.
}

func (a Any) GetRef() int {
	return AnyRef
}
