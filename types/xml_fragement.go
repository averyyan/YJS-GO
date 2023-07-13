package types

import (
	"YJS-GO/structs"
	"YJS-GO/utils"
)

type XmlFragment struct {
	AbstractType
	prelimContent []any
}

func (x XmlFragment) firstChild() any {
	first := x.first
	if first != nil {
		return first.Content.GetContent()[0]
	}
	return nil
}

func (x XmlFragment) integrate(doc *utils.YDoc, item *structs.Item) {
	x.Doc = doc
	x.Item = item
	x.insert(0, x.prelimContent)
	x.prelimContent = nil
}

func (x XmlFragment) insert(i int, content any) {
	if x.Doc != nil {
		transact(this.doc, transaction => {
			typeListInsertGenerics(transaction, this, index, content)
		})
	} else {
		// @ts-ignore _prelimContent is defined because this is not yet integrated
		x.prelimContent.splice(index, 0, ...content)
	}
}
