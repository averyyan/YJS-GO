package content

import (
	"encoding/json"

	"YJS-GO/structs"
	"YJS-GO/utils"
)

var _ structs.IContentExt = (*Json)(nil)

type Json struct {
	Content []any
}

func NewJson(data []any) *Json {
	return &Json{Content: data}
}

func (j Json) Integrate(transaction *utils.Transaction, item *structs.Item) {
	// Do nothing.

}

func (j Json) SetRef(i int) {
	// Do nothing.
}

func ReadJson(decoder utils.IUpdateDecoder) (*Json, error) {
	var length = decoder.ReadLength()
	var content = make([]any, length)

	for i := 0; i < int(length); i++ {
		var jsonStr = decoder.ReadString()
		var jsonObj interface{}
		if jsonStr == "undefined" {
			jsonObj = nil
		} else {
			json.Unmarshal([]byte(jsonStr), jsonObj)
		}
		content = append(content, jsonObj)
	}

	return NewJson(content), nil
}

func (j Json) Copy() structs.IContent {
	return NewJson(j.Content)
}

func (j Json) Splice(offset uint64) structs.IContent {
	var right = NewJson(j.Content[int(offset) : len(j.Content)-int(offset)])
	j.Content = append(j.Content[0:offset], j.Content[len(j.Content)-int(offset):]...)
	return right
}

func (j Json) MergeWith(right structs.IContent) bool {
	// Debug.Assert(right is ContentJson)
	j.Content = append(j.Content, right.GetContent().([]any)...)
	return true
}

func (j Json) GetContent() any {
	return j.Content
}

func (j Json) GetLength() int {
	if j.Content != nil {
		return len(j.Content)
	}
	return 0
}

func (j Json) Countable() bool {
	return true
}

func (j Json) Write(encoder utils.IUpdateEncoder, offset int) {
	var length = len(j.Content)
	encoder.WriteLength(length)
	for i := offset; i < length; i++ {
		jsonStr, err := json.Marshal(j.Content[i])
		if err != nil {
			continue
		}
		encoder.WriteString(string(jsonStr))
	}
}

func (j Json) Gc(store *utils.StructStore) {
	// Do nothing.

}

func (j Json) Delete(transaction *utils.Transaction) {
	// Do nothing.

}

func (j Json) GetRef() int {
	return JsonRef
}
