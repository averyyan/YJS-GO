package content

import (
	"container/list"

	"YJS-GO/structs"
	"YJS-GO/utils"
)

var _ structs.IContentExt = (*Embed)(nil)

type Embed struct {
}

func ReadEmbed(decoder utils.IUpdateDecoder) (Embed, error) {
	// TODO implement me
	panic("implement me")
}

func (e Embed) Copy() structs.IContent {
	// TODO implement me
	panic("implement me")
}

func (e Embed) Splice(offset uint64) structs.IContent {
	// TODO implement me
	panic("implement me")
}

func (e Embed) MergeWith(right structs.IContent) bool {
	// TODO implement me
	panic("implement me")
}

func (e Embed) GetContent() list.List {
	// TODO implement me
	panic("implement me")
}

func (e Embed) GetLength() int {
	// TODO implement me
	panic("implement me")
}

func (e Embed) Countable() bool {
	// TODO implement me
	panic("implement me")
}

func (e Embed) Write(encoder utils.IUpdateEncoder, offset int) {
	// TODO implement me
	panic("implement me")
}

func (e Embed) Gc(store *utils.StructStore) {
	// TODO implement me
	panic("implement me")
}

func (e Embed) Delete(transaction *utils.Transaction) {
	// TODO implement me
	panic("implement me")
}

func (e Embed) Integrate(transaction *utils.Transaction, item structs.Item) {
	// TODO implement me
	panic("implement me")
}

func (e Embed) GetRef() int {
	// TODO implement me
	panic("implement me")
}
