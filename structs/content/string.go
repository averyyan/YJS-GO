package content

import (
	"strings"

	"YJS-GO/structs"
	"YJS-GO/utils"
)

var _ structs.IContentExt = (*String)(nil)

type String struct {
	content []any
}

func NewString(value string) *String {
	return &String{}
}

func (s String) SetRef(i int) {
	// Do nothing.
}

func ReadString(decoder utils.IUpdateDecoder) (*String, error) {
	// TODO implement me
	panic("implement me")
}

func (s String) Copy() structs.IContent {
	// TODO implement me
	panic("implement me")
}

func (s String) Splice(offset uint64) structs.IContent {
	// TODO implement me
	panic("implement me")
}

func (s String) MergeWith(right structs.IContent) bool {
	// TODO implement me
	panic("implement me")
}

func (s String) GetContent() any {
	return s.content
}

func (s String) GetLength() int {
	// TODO implement me
	panic("implement me")
}

func (s String) Countable() bool {
	// TODO implement me
	panic("implement me")
}

func (s String) Write(encoder utils.IUpdateEncoder, offset int) {
	// TODO implement me
	panic("implement me")
}

func (s String) Gc(store *utils.StructStore) {
	// TODO implement me
	panic("implement me")
}

func (s String) Delete(transaction *utils.Transaction) {
	// TODO implement me
	panic("implement me")
}

func (s String) Integrate(transaction *utils.Transaction, item *structs.Item) {
	// TODO implement me
	panic("implement me")
}

func (s String) GetRef() int {
	return StringRef
}

func (s String) GetString() string {
	var a = strings.Builder{}
	for _, str := range s.GetContent().([]any) {
		a.WriteString(str.(string))
	}
	return a.String()
}
func (s String) AppendToBuilder(sb strings.Builder) {
	for _, c := range s.content {
		sb.Append((char)
		c)
	}
}
