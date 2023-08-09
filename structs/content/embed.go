package content

import (
	"YJS-GO/structs"
	"YJS-GO/utils"
)

var _ structs.IContentExt = (*Embed)(nil)

type Embed struct {
	Embed any
}

func NewEmbed(embed any) *Embed {
	var ret = &Embed{}
	ret.Embed = embed
	return ret
}

func (e Embed) SetRef(i int) {
	// Do nothing.
}

func ReadEmbed(decoder utils.IUpdateDecoder) (*Embed, error) {
	var content = decoder.ReadJson()
	return NewEmbed(content), nil
}

func (e Embed) Copy() structs.IContentExt {
	return NewEmbed(e.Embed)
}

func (e Embed) Splice(offset uint64) structs.IContentExt {
	// Do nothing.
	return nil
}

func (e Embed) MergeWith(right structs.IContentExt) bool {
	return false
}

func (e Embed) GetContent() any {
	return []any{e.Embed}
}

func (e Embed) GetLength() int {
	return 1
}

func (e Embed) Countable() bool {
	return true
}

func (e Embed) Write(encoder utils.IUpdateEncoder, offset int) {
	encoder.WriteJson(e.Embed)
}

func (e Embed) Gc(store *utils.StructStore) {
	// Do nothing.
}

func (e Embed) Delete(transaction *utils.Transaction) {
	// Do nothing.
}

func (e Embed) Integrate(transaction *utils.Transaction, item *structs.Item) {
	// Do nothing.
}

func (e Embed) GetRef() int {
	return EmbedRef
}
