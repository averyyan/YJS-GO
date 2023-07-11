package types

import (
	"YJS-GO/structs"
	"YJS-GO/utils"
)

type EventHandler func(v any)

type AbstractType struct {
	Item    *structs.Item
	Start   *structs.Item
	ItemMap map[string]*structs.Item

	EventHandler     EventHandler
	DeepEventHandler EventHandler

	Doc    *utils.YDoc
	Parent *AbstractType
	Length uint64
}

func (at AbstractType) Integrate(doc *utils.YDoc, item *structs.Item) {
	at.Doc = doc
	at.Item = item
}

type YEventArgs struct {
	Event       *utils.YEvent
	Transaction *utils.Transaction
}

type YDeepEventArgs struct {
	Events      []*utils.YEvent
	Transaction *utils.Transaction
}
