package types

import (
	"YJS-GO/structs"
	"YJS-GO/utils"
)

type XmlElement struct {
	XmlFragment
	NodeName    string
	PrelimAttrs map[string]string
}

func NewXmlElement(name string) *XmlElement {
	t := &XmlElement{XmlFragment: XmlFragment{
		AbstractType:  AbstractType{},
		prelimContent: nil,
	}, NodeName: name, PrelimAttrs: map[string]string{}}
	return t
}

func (e XmlElement) nextSibling() XmlFragment {
	var t *structs.Item
	if e.Item != nil {
		t = e.Item.Next.(*structs.Item)
	}
	if t != nil {
		return t.Content
	}
	return nil
}

func ReadXmlElement(decoder utils.IUpdateDecoder) *XmlElement {
	return NewXmlElement(decoder.ReadKey())
}
