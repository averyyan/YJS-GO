package structs

import (
	"reflect"

	"YJS-GO/lib0"
	"YJS-GO/structs/content"
	"YJS-GO/types"
	"YJS-GO/utils"
)

const (
	Keep      InfoEnum = 1 << 0
	Countable InfoEnum = 1 << 1
	Deleted   InfoEnum = 1 << 2
	Marker    InfoEnum = 1 << 3
)

type InfoEnum int

var _ IAbstractStruct = (*Item)(nil)

type Item struct {
	Id          *utils.ID
	Length      uint64
	Deleted     bool
	Info        InfoEnum
	LeftOrigin  *utils.ID
	Left        any // AbstractStruct
	RightOrigin *utils.ID
	Right       any // AbstractStruct
	Parent      any // Object
	// If the parent refers to this item with some kind of key (e.g. YMap).
	// The key is then used to refer to the list in which to insert this item.
	// If 'parentSub =nil', type._start is the list in which to insert to.
	// Otherwise, it is 'parent._map'.
	ParentSub string
	Redone    *utils.ID
	Content   IContentExt
	Marker    bool
	Keep      bool
	Countable bool
	LastId    *utils.ID
	Next      any // AbstractStruct
	Prev      any // AbstractStruct
}

func (i *Item) GetLength() uint64 {
	return i.Length
}

func (i *Item) GetDeleted() bool {
	return i.Deleted
}

func (i *Item) ID() *utils.ID {
	return i.Id
}

// MergeWith func (i *Item) MergeWith(right AbstractStruct) bool {
func (i *Item) MergeWith(right any) bool {
	var rightItem = i.Right.(*Item)
	if utils.EQ(rightItem.LeftOrigin, i.LastId) &&
		i.Right == right &&
		utils.EQ(rightItem.RightOrigin, i.RightOrigin) &&
		i.Id.Client == right.(*Item).Id.Client &&
		i.Id.Clock+i.Length == right.(*Item).Id.Clock &&
		i.Deleted == right.(*Item).Deleted &&
		i.Redone == nil &&
		reflect.TypeOf(i.Content) == reflect.TypeOf(rightItem.Content) &&
		i.Content.MergeWith(rightItem.Content) {
		if rightItem.Keep {
			i.Keep = true
		}
		Right := rightItem.Right.(*Item)
		if reflect.TypeOf(Right) == reflect.TypeOf(&Item{}) {
			Right.Left = i
		}
		i.Length += rightItem.Length
		return true
	}
	return false
}

func (i *Item) Delete(transaction *utils.Transaction) {
	if !i.Deleted {
		var parent = i.Parent
		if i.Countable && i.ParentSub == "" {
			if parent != nil {
				parent.(*types.AbstractType).Length -= i.Length
			}
		}
		i.MarkDeleted()
		transaction.DeleteSet.Add(i.Id.Client, i.Id.Clock, uint64(i.Length))
		transaction.AddChangedTypeToTransaction(parent.(*types.AbstractType), i.ParentSub)
		i.Content.Delete(transaction)
	}
}

func (i *Item) Integrate(transaction *utils.Transaction, offset int) {
	if offset > 0 {
		i.Id = &utils.ID{Client: i.Id.Client, Clock: i.Id.Clock + uint64(offset)}
		i.Left = transaction.Doc.Store.GetItemCleanEnd(transaction, &utils.ID{Client: i.Id.Client, Clock: i.Id.Clock - 1})
		i.LeftOrigin = i.Left.(*Item).LastId
		i.Content = i.Content.Splice(uint64(offset))
		i.Length -= uint64(offset)
	}

	if i.Parent != nil {
		if (i.Left == nil && (i.Right == nil || i.Right.(*Item).Left != nil)) ||
			i.Left != nil && i.Left.(*Item).Right != i.Right {
			var left = i.Left.(*Item)
			var o IAbstractStruct
			// Set 'o' to the first conflicting item.
			if left != nil {
				o = left.Right.(*Item)
			} else if i.ParentSub != "" {
				// Debug.Assert(Parent is AbstractType);
				item, ok := i.Parent.(types.AbstractType).ItemMap[i.ParentSub]
				if ok {
					o = item
				}
				for o != nil && o.(*Item).Left != nil {
					o = o.(*Item).Left.(*Item)
				}
			} else {
				// Debug.Assert(ParentisAbstractType)o = Parent.(AbstractType)?._start
			}

			var conflictingItems = map[IAbstractStruct]struct{}{}
			var itemsBeforeOrigin = map[IAbstractStruct]struct{}{}

			for o != nil && o != i.Right {
				itemsBeforeOrigin[o] = struct{}{}
				conflictingItems[o] = struct{}{}

				if utils.EQ(i.LeftOrigin, o.(*Item).LeftOrigin) {
					// Case 1
					if o.ID().Client < i.ID().Client {
						left = o.(*Item)
						conflictingItems = map[IAbstractStruct]struct{}{}
					} else if utils.EQ(i.RightOrigin, o.(*Item).RightOrigin) {
						// This and 'o' are conflicting and point to the same integration points.
						// The id decides which item comes first.
						// Since this is to the left of 'o', we can break here.
						break
					}
					// Else, 'o' might be integrated before an item that this conflicts with.
					// If so, we will find it in the next iterations.
					// Use 'Find' instead of 'GetItemCleanEnd', because we don't want / need to split items.
				} else {
					leftValue := transaction.Doc.Store.Find(o.(*Item).LeftOrigin)
					_, ok := itemsBeforeOrigin[leftValue]
					if o.(*Item).LeftOrigin != nil && ok {
						// Case 2
						// TODO: Store.Find is called twice here, call once?
						_, ok := conflictingItems[leftValue]
						if !ok {
							left = o.(*Item)
							conflictingItems = map[IAbstractStruct]struct{}{}
						}
					} else {
						break
					}
				}
				o = o.(*Item).Right.(*Item)
			}
			i.Left = left
		}

		// Reconnect left/right + update parent map/start if necessary.
		if i.Left != nil {
			leftItem, ok := i.Left.(*Item)
			if ok {
				var right = leftItem.Right
				i.Right = right
				leftItem.Right = i
			} else {
				i.Right = nil
			}
		} else {
			var r IAbstractStruct
			if i.ParentSub != "" {
				var item *Item
				item, ok := i.Parent.(*types.AbstractType).ItemMap[i.ParentSub]
				if ok {
					r = item
				}
				for r != nil && r.(*Item).Left != nil {
					r = r.(*Item).Left.(*Item)
				}
			} else {
				abstractTypeParent, ok := i.Parent.(*types.AbstractType)
				if ok {
					r = abstractTypeParent.Start
					abstractTypeParent.Start = i
				} else {
					r = nil
				}
			}
			i.Right = r
		}

		if i.Right != nil {
			rightItem, ok := i.Right.(*Item)
			if ok {
				rightItem.Left = i
			}
		} else if i.ParentSub != "" {
			// Set as current parent value if right == nil and this is parentSub.
			i.Parent.(*types.AbstractType).ItemMap[i.ParentSub] = i
			// This is the current attribute value of parent. Delete right.
			i.Left.(IAbstractStruct).Delete(transaction)
		}

		// Adjust length of parent.
		if i.ParentSub == "" && i.Countable && !i.Deleted {
			// Debug.Assert(Parent is AbstractType)
			i.Parent.(*types.AbstractType).Length += i.Length
		}

		transaction.Doc.Store.AddStruct(i)
		i.Content.(IContentExt).Integrate(transaction, i)

		// Add parent to transaction.changed.
		transaction.AddChangedTypeToTransaction(i.Parent.(*types.AbstractType), i.ParentSub)

		if (i.Parent.(*types.AbstractType).Item != nil && i.Parent.(*types.AbstractType).Item.Deleted) || (i.ParentSub != "" && i.Right != nil) {
			// Delete if parent is deleted or if this is not the current attribute value of parent.
			i.Delete(transaction)
		}
	} else {
		// Parent is not defined. Integrate GC struct instead.
		(&GC{Id: i.Id, Length: i.Length}).Integrate(transaction, 0)
	}
}

func (i *Item) GetMissing(transaction *utils.Transaction, store *utils.StructStore) (uint64, error) {
	// TODO implement me
	panic("implement me")
}

func (i *Item) Write(encoder utils.IUpdateEncoder, offset int) {
	var origin = &utils.ID{Client: i.Id.Clock, Clock: i.Id.Clock + uint64(offset) - 1}
	if offset <= 0 {
		origin = i.LeftOrigin
	}
	var rightOrigin = i.RightOrigin
	var parentSub = i.ParentSub
	b := uint(0)
	if origin != nil {
		b = lib0.Bit8
	}
	c := uint(0)
	if rightOrigin != nil {
		c = lib0.Bit7
	}
	d := uint(0)
	if parentSub != "" {
		d = lib0.Bit6
	}
	var info = uint(i.Content.GetRef())&lib0.Bits5 | b | c | d

	encoder.WriteInfo(uint8(info))

	if origin != nil {
		encoder.WriteLeftId(origin)
	}

	if rightOrigin != nil {
		encoder.WriteRightId(rightOrigin)
	}

	if origin == nil && rightOrigin == nil {
		var parent = i.Parent
		var parentItem = parent.(*types.AbstractType).Item
		if parentItem == nil {
			// parent type on y._map.
			// find the correct key
			var yKey = parent.(*types.AbstractType).FindRootTypeKey()
			encoder.WriteParentInfo(true)
			encoder.WriteString(yKey)
		} else {
			encoder.WriteParentInfo(false)
			encoder.WriteLeftId(parentItem.Id)
		}

		if parentSub != "" {
			encoder.WriteString(parentSub)
		}
	}

	i.Content.Write(encoder, offset)
}

func (i *Item) SetMarker(value bool) {
	if value {
		i.Info |= Marker
	} else {
		i.Info &= ^Marker
	}
}

func (i *Item) SetKeep(value bool) {
	if value {
		i.Info |= Keep
	} else {
		i.Info &= ^Keep
	}
}

func (i *Item) SetCountable(value bool) {
	if value {
		i.Info |= Countable
	} else {
		i.Info &= ^Countable
	}
}

func (i *Item) GetNext() any {
	var n = i.Right
	for n != nil && n.(*Item).Deleted {
		b := n
		n = b.(*Item).Next
	}
	return n
}

func (i *Item) GetPrev() any {
	var n = i.Left
	for n != nil && n.(*Item).Deleted {
		b := n
		n = b.(*Item).Left
	}
	return n
}

func (i *Item) MarkDeleted() {
	i.Info |= Deleted
}

// SplitItem Split 'leftItem' into two items.
func (i *Item) SplitItem(transaction *utils.Transaction, diff uint64) *Item {
	var client = i.Id.Client
	var clock = i.Id.Clock
	var countable InfoEnum
	if i.Content.Countable() {
		countable = Countable
	}
	var rightItem = &Item{
		Id: &utils.ID{
			Client: client,
			Clock:  clock + diff,
		},
		LeftOrigin: &utils.ID{
			Client: client,
			Clock:  clock + diff - 1,
		},
		Left:        i,
		RightOrigin: i.RightOrigin,
		Right:       i.Right,
		Parent:      i.Parent,
		ParentSub:   i.ParentSub,
		Content:     i.Content.Splice(diff).(IContentExt),
		Info:        countable,
	}

	if i.Deleted {
		rightItem.MarkDeleted()
	}
	if i.Keep {
		rightItem.Keep = true
	}
	if i.Redone != nil {
		rightItem.Redone = &utils.ID{
			Client: i.Redone.Client,
			Clock:  i.Redone.Clock + diff,
		}
	}

	// Update left (do not set leftItem.RightOrigin as it will lead to problems when syncing).
	i.Right = rightItem
	// Update right.
	rightIt := rightItem.Right.(*Item)
	rightIt.Left = rightItem
	// Right is more specific.
	transaction.MergeStructs = append(transaction.MergeStructs, rightItem)
	// Update parent._map.
	if rightItem.ParentSub != "" && rightItem.Right == nil {
		rightItem.Parent.(types.AbstractType).ItemMap[rightItem.ParentSub] = rightItem
	}

	i.Length = diff
	return rightItem
}

func (i *Item) Gc(store *utils.StructStore,
	parentGCd bool) {
	if !i.Deleted {
		return
		// throw new InvalidOperationException();
	}
	i.Content.(IContentExt).Gc(store)
	if parentGCd {
		store.ReplaceStruct(i, &GC{Id: i.Id, Length: i.Length})
	} else {
		i.Content = &content.Deleted{Length: i.Length}
	}
}

func (i *Item) KeepItemAndParents(value bool) {
	var item = i
	for item != nil && item.Keep != value {
		item.Keep = value
		tmp, ok := item.Parent.(*types.AbstractType)
		if ok && tmp != nil {
			item = tmp.Item
		} else {
			item = nil
		}
	}
}

func NewItem(id *utils.ID, left IAbstractStruct, leftOrigin *utils.ID, right IAbstractStruct, rightOrigin *utils.ID,
	parent any, parentSub string, content IContent) *Item {
	t := &Item{}
	t.Id = id
	t.Length = uint64(content.GetLength())
	t.LeftOrigin = leftOrigin
	t.Left = left
	t.Right = right
	t.RightOrigin = rightOrigin
	t.Parent = parent
	t.ParentSub = parentSub
	t.Redone = nil
	t.Content = content.(IContentExt)
	if content.Countable() {
		t.Info = Countable
	} else {
		t.Info = 0
	}
	return t
}
