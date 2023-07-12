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

// CallObserver Creates YEvent and calls all type observers.
// Must be implemented by each type.
func (at AbstractType) CallObserver(transaction *utils.Transaction, subs map[string]struct{}) {
	// Do nothing.
}

func (at AbstractType) CallDeepEventHandlerListeners(events []*utils.YEvent, transaction *utils.Transaction) {
	if at.DeepEventHandler != nil {
		at.DeepEventHandler(&YDeepEventArgs{
			Events:      events,
			Transaction: transaction,
		})
	}
}

type YEventArgs struct {
	Event       *utils.YEvent
	Transaction *utils.Transaction
}

type YDeepEventArgs struct {
	Events      []*utils.YEvent
	Transaction *utils.Transaction
}
