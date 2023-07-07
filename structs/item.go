package structs

import (
	"reflect"

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
	Id          utils.ID
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
	// If 'parentSub = null', type._start is the list in which to insert to.
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

func (i *Item) GetLength() uint {
	return uint(i.Length)
}

func (i *Item) GetDeleted() bool {
	return i.Deleted
}

func (i *Item) ID() utils.ID {
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

func (i *Item) Delete(transaction utils.Transaction) {
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

func (i *Item) Integrate(transaction utils.Transaction, offset int) {
	// TODO implement me
	panic("implement me")
}

func (i *Item) GetMissing(transaction utils.Transaction, store utils.StructStore) {
	// TODO implement me
	panic("implement me")
}

func (i *Item) Write(encoder utils.IUpdateEncoder, offset int) {
	// TODO implement me
	panic("implement me")
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
