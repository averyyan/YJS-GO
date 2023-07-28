package content

import (
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
	var len = decoder.ReadLength();
	var content = new List<object>(len);

	for (int i = 0; i < len; i++)
	{
	var jsonStr = decoder.ReadString();
	object jsonObj = string.Equals(jsonStr, "undefined")
	? null
	: Newtonsoft.Json.JsonConvert.DeserializeObject(jsonStr);
	content.Add(jsonObj);
	}

	return new ContentJson(content);
}

func (j Json) Copy() structs.IContent {
	return NewJson(j.Content)
}

func (j Json) Splice(offset uint64) structs.IContent {
	var right = NewJson(j.Content[int(offset) : len(j.Content)-int(offset)])
	j.Content.RemoveRange(offset, _content.Count-offset)
	return right
}

func (j Json) MergeWith(right structs.IContent) bool {
	Debug.Assert(right is ContentJson);
	_content.AddRange((right as ContentJson)._content);
	return true;
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
	// TODO implement me
	panic("implement me")
}

func (j Json) Write(encoder utils.IUpdateEncoder, offset int) {
	var len = _content.Count;
	encoder.WriteLength(len);
	for (int i = offset; i < len; i++)
	{
	var jsonStr = Newtonsoft.Json.JsonConvert.SerializeObject(_content[i]);
	encoder.WriteString(jsonStr);
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
