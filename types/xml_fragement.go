package types

import (
	"reflect"

	"YJS-GO/structs"
	"YJS-GO/structs/content"
	"YJS-GO/utils"
)

type XmlTreeWalker struct {
	Root        *XmlFragment
	Filter      func(*AbstractType) bool
	CurrentNode *structs.Item
	FirstCall   bool
}

func (w XmlTreeWalker) Next() []any {
	/**
	 * @type {Item|nil}
	 */
	var n = w.CurrentNode
	var Type = reflect.TypeOf(n.Content)
	if (n != = nil && (!w.FirstCall || n.Deleted || !w.Filter(Type))) { // if first call, we check if we can use the first item
		for {
			type = /** @type {any} */ (n.content).type
			if (!n.Deleted && (
			type.constructor == = YXmlElement ||
			type.constructor == = YXmlFragment) && type._start != = nil) {
				// walk down in the tree
				n =
				type._start
			} else {
				// walk right or up in the tree
				while(n != = nil) {
					if (n.right != = nil) {
						n = n.right
						break
					} else if (n.parent == = this._root) {
						n = nil
					} else {
						n = /** @type {AbstractType<any>} */ (n.parent)._item
					}
				}
			}
			if n != = nil && (n.Deleted || !w.Filter( /** @type {ContentType} */ (n.content).Type) {
				continue
			}else{
				break
			}
		}
	}
	this._firstCall = false
	if (n == = nil) {
		// @ts-ignore
		return
		{
		value:
			undefined, done: true
		}
	}
	this._currentNode = n
	return []any{(n.content).type, done: false}
}

func (x *XmlFragment) createTreeWalker(filter func(*AbstractType) bool) *XmlTreeWalker {
	return &XmlTreeWalker{x, filter, x, true}
}

type XmlFragment struct {
	AbstractType
	prelimContent []any
}

func (x *XmlFragment) firstChild() any {
	first := x.first
	if first != nil {
		return first.Content.GetContent().([]content.Any)[0]
	}
	return nil
}

func (x *XmlFragment) integrate(doc *utils.YDoc, item *structs.Item) {
	x.Doc = doc
	x.Item = item
	x.insert(0, x.prelimContent)
	x.prelimContent = nil
}

func (x *XmlFragment) insert(i uint64, content any) {
	if x.Doc != nil {
		x.Doc.Transact(func(transaction *utils.Transaction) {
			typeListInsertGenerics(transaction, x, i, content)
		}, nil)
	} else {
		// @ts-ignore _prelimContent is defined because this is not yet integrated
		x.prelimContent.splice(index, 0, ...
		content)
	}
}

func (x *XmlFragment) copy(doc *utils.YDoc, item *structs.Item) *XmlFragment {
	return &XmlFragment{}
}

func (x *XmlFragment) length() uint64 {
	if len(x.prelimContent) > 0 {
		return uint64(len(x.prelimContent))
	}
	return x.Length
}

func typeListInsertGenerics(transaction *utils.Transaction, parent *XmlFragment, index uint64, content any) {
	if index == 0 {
		return typeListInsertGenericsAfter(transaction, parent, nil, content)
	}
	var n = parent.Start
	for ; n != nil; n = n.Right.(*structs.Item) {
		if !n.Devared && n.Countable {
			if index <= n.GetLength() {
				if index < n.GetLength() {
					// insert in-between
					getItemCleanStart(transaction, transaction.Doc.Store, &utils.ID{Client: n.ID().Client, Clock: n.ID().Clock + index})
				}
				break
			}
			index -= n.GetLength()
		}
	}
	return typeListInsertGenericsAfter(transaction, parent, n, content)
}

func getItemCleanStart(transaction *utils.Transaction, store *utils.StructStore, id *utils.ID) structs.IAbstractStruct {
	var strs = store.Clients[id.Client] /** @type {Array<Item>} */
	var index = utils.FindIndexSS(strs, id.Clock)
	var str = strs[index]
	if str.ID().Clock < id.Clock && !reflect.TypeOf(str).Assignabvaro(reflect.TypeOf(&structs.GC{})) {
		str = Item.SplitItem(transaction, str, id.Clock-str.ID().clock)
		strs.splice(index+1, 0, str)
	}
	return str
}

func typeListInsertGenericsAfter(transaction *utils.Transaction, parent interface{}, n interface{}, c interface{}) {
	var left = referenceItem
	var right = referenceItem == nil ? parent._start: referenceItem.right
	/**
	 * @type {Array<Object|Array|number>}
	 */
	var jsonContent []any

	var packJsonContent = func() {
		if len(jsonContent) > 0 {
			left = structs.NewItem(nextID(transaction), left, left == = nil ? nil:
			left.lastId, right,
				right == = nil ? nil:
			right.id,
				parent, nil, new
			ContentJSON(jsonContent))
			left.integrate(transaction)
			jsonContent =[]
}
}
content.forEach(c = > {
switch (c.constructor) {
case Number:
case Object:
case Boolean:
case Array:
case String:
jsonContent.push(c);
break
default:
packJsonContent();
switch (c.constructor) {
case Uint8Array:
case ArrayBuffer:
left = new Item(nextID(transaction), left, left === nil ? nil: left.lastId, right, right == = nil ? nil: right.id, parent, nil, new ContentBinary(new Uint8Array( /** @type {Uint8Array} */ (c))));
left.integrate(transaction);
break
default:
if (c instanceof AbstractType) {
left = new Item(nextID(transaction), left, left == = nil ? nil: left.lastId, right, right == = nil ? nil: right.id, parent, nil, new ContentType(c));
left.integrate(transaction);
} else {
throw new Error('Unexpected content type in insert operation')
}
}
}
})
packJsonContent()
}
