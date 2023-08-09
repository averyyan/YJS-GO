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
	runeArray := []rune(value)
	var tArr []any
	for _, r := range runeArray {
		tArr = append(tArr, r)
	}
	return &String{tArr}
}

func (s String) SetRef(i int) {
	// Do nothing.
}

func ReadString(decoder utils.IUpdateDecoder) (*String, error) {
	return NewString(decoder.ReadString()), nil
}

func (s String) Copy() structs.IContentExt {
	return NewString(s.GetString())
}

func (s String) Splice(offset uint64) structs.IContentExt {
	var t = s.content[int(offset) : len(s.content)-int(offset)]
	var sb = &strings.Builder{}
	for i := 0; i < len(t); i++ {
		sb.WriteRune(t[i].(rune))
	}
	var right = NewString(sb.String())
	s.content = append(s.content[:offset], s.content[len(s.content)-int(offset):])

	// Prevent encoding invalid documents because of splitting of surrogate pairs.
	var firstCharCode = s.content[offset-1].(rune)
	if firstCharCode >= 0xD800 && firstCharCode <= 0xDBFF {
		// Last character of the left split is the start of a surrogate utf16/ucs2 pair.
		// We don't support splitting of surrogate pairs because this may lead to invalid documents.
		// Replace the invalid character with a unicode replacement character U+FFFD.
		s.content[offset-1] = '\uFFFD'

		// Replace right as well.
		right.content[0] = '\uFFFD'
	}

	return right
}

func (s String) MergeWith(right structs.IContentExt) bool {
	// Debug.Assert(right is ContentString);
	s.content = append(s.content, (right.(String)).content)
	return true
}

func (s String) GetContent() any {
	return s.content
}

func (s String) GetLength() int {
	return len(s.content)
}

func (s String) Countable() bool {
	return true
}

func (s String) Write(encoder utils.IUpdateEncoder, offset int) {
	// var sb = new StringBuilder(_content.Count - offset);
	var sb = &strings.Builder{}
	for i := offset; i < len(s.content); i++ {
		sb.WriteRune(s.content[i].(rune))
	}
	var str = sb.String()
	encoder.WriteString(str)
}

func (s String) Gc(store *utils.StructStore) {
	// Do nothing.
}

func (s String) Delete(transaction *utils.Transaction) {
	// Do nothing.
}

func (s String) Integrate(transaction *utils.Transaction, item *structs.Item) {
	// Do nothing.
}

func (s String) GetRef() int {
	return StringRef
}

func (s String) GetString() string {
	var a = &strings.Builder{}
	for _, str := range s.GetContent().([]any) {
		a.WriteString(str.(string))
	}
	return a.String()
}
func (s String) AppendToBuilder(sb *strings.Builder) {
	for _, c := range s.content {
		sb.WriteString(c.(string))
	}
}
