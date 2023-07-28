package content

import (
	"YJS-GO/structs"
	"YJS-GO/utils"
)

var _ structs.IContentExt = (*Doc)(nil)

type Doc struct {
	Doc  *utils.YDoc
	Opts *utils.YDocOptions
}

func NewDoc(doc *utils.YDoc) *Doc {
	Opts := utils.YDocOptions{}

	if !doc.GC {
		Opts.GC = false
	}

	if doc.AutoLoad {
		Opts.AutoLoad = true
	}

	if doc.Meta != nil {
		Opts.Meta = doc.Meta
	}
	return &Doc{Doc: doc}
}

func (d Doc) SetRef(i int) {
	// Do nothing.
}

func ReadDoc(decoder utils.IUpdateDecoder) (*Doc, error) {
	var guidStr = decoder.ReadString()
	var opts = utils.ReadYDocOptions(decoder)
	opts.GUID = guidStr
	return NewDoc(utils.NewDoc(opts)), nil
}

func (d Doc) Copy() structs.IContent {
	return NewDoc(d.Doc)
}

func (d Doc) Splice(offset uint64) structs.IContent {
	return nil
}

func (d Doc) MergeWith(right structs.IContent) bool {
	return false
}

func (d Doc) GetContent() any {
	return []any{NewDoc(d.Doc)}
}

func (d Doc) GetLength() int {
	return 1
}

func (d Doc) Countable() bool {
	return true
}

func (d Doc) Write(encoder utils.IUpdateEncoder, offset int) {
	// 32 digits separated by hyphens, no braces.
	encoder.WriteString(d.Doc.GUID)
	d.Opts.Write(encoder, offset)
}

func (d Doc) Gc(store *utils.StructStore) {
	// Do nothing.
}

func (d Doc) Delete(transaction *utils.Transaction) {
	if _, ok := transaction.SubdocsAdded[d.Doc]; ok {
		delete(transaction.SubdocsAdded, d.Doc)
	} else {
		transaction.SubdocsRemoved[d.Doc] = struct{}{}
	}
}

func (d Doc) Integrate(transaction *utils.Transaction, item *structs.Item) {
	// This needs to be reflected in doc.destroy as well.
	d.Doc.Item = item
	transaction.SubdocsAdded[d.Doc] = struct{}{}

	if d.Doc.ShouldLoad {
		transaction.SubdocsLoaded[d.Doc] = struct{}{}
	}
}

func (d Doc) GetRef() int {
	return DocRef
}
