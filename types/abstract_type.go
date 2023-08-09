package types

import (
	"fmt"

	"YJS-GO/structs"
	"YJS-GO/utils"
)

type EventHandler func(v any)

type IAbstractType interface {
	CallObserver(transaction *utils.Transaction, subs map[string]struct{})
	CallDeepEventHandlerListeners(events []*utils.YEvent, transaction *utils.Transaction)
	Integrate(doc *utils.YDoc, item *structs.Item)
}

type AbstractType struct {
	Item    *structs.Item
	Start   *structs.Item
	ItemMap map[string]*structs.Item

	EventHandler     EventHandler
	DeepEventHandler EventHandler

	Doc    *utils.YDoc
	Parent *AbstractType
	Length uint64
	first  *structs.Item
}

func (at *AbstractType) Integrate(doc *utils.YDoc, item *structs.Item) {
	at.Doc = doc
	at.Item = item
}

// CallObserver Creates YEvent and calls all type observers.
// Must be implemented by each type.
func (at *AbstractType) CallObserver(transaction *utils.Transaction, subs map[string]struct{}) {
	// Do nothing.
}

func (at *AbstractType) CallDeepEventHandlerListeners(events []*utils.YEvent, transaction *utils.Transaction) {
	if at.DeepEventHandler != nil {
		at.DeepEventHandler(&YDeepEventArgs{
			Events:      events,
			Transaction: transaction,
		})
	}
}

func (at *AbstractType) Write(utils.IUpdateEncoder) {
	fmt.Printf("not implement")
}

func (at *AbstractType) FindRootTypeKey() string {
	return at.Doc.FindRootTypeKey(at)
}

type YEventArgs struct {
	Event       *utils.YEvent
	Transaction *utils.Transaction
}

type YDeepEventArgs struct {
	Events      []*utils.YEvent
	Transaction *utils.Transaction
}
