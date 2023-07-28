package types

import (
	"math"

	"YJS-GO/structs"
	"YJS-GO/structs/content"
	"YJS-GO/utils"
)

type YArrayBase struct {
	AbstractType
	searchMarkers *ArraySearchMarkerCollection
}

func NewYArrayBase() *YArrayBase {
	return &YArrayBase{
		AbstractType:  AbstractType{},
		searchMarkers: &ArraySearchMarkerCollection{},
	}
}

func (b YArrayBase) ClearSearchMarkers() {
	b.searchMarkers.Clear()
}

func (b YArrayBase) InsertGeneric(transaction *utils.Transaction, index uint64, content []any) {
	if index == 0 {
		if b.searchMarkers.Count() > 0 {
			b.searchMarkers.UpdateMarkerChanges(index, uint64(len(content)))
		}
		b.InsertGenericsAfter(transaction, nil, content)
		return
	}

	startIndex := index
	var marker = b.FindMarker(index)
	var n = b.Start

	if marker != nil {
		n = marker.P
		index -= marker.Index

		// We need to iterate one to the left so that the algorithm works.
		if index == 0 {
			// @todo: refactor this as it actually doesn't consider formats.
			n = n.Prev.(*structs.Item)

			if n != nil && n.Countable && !n.Deleted {
				index += n.GetLength()
			}
		}
	}

	for ; n != nil; n = n.Right.(*structs.Item) {
		if !n.Deleted && n.Countable {
			if index <= n.Length {
				if index < n.Length {
					// insert in-between
					transaction.Doc.Store.GetItemCleanStart(transaction, &utils.ID{Client: n.Id.Client, Clock: n.Id.Clock + uint64(index)})
				}
				break
			}

			index -= n.Length
		}
	}

	if b.searchMarkers.Count() > 0 {
		b.searchMarkers.UpdateMarkerChanges(startIndex, uint64(len(content)))
	}

	b.InsertGenericsAfter(transaction, n, content)
}

// CallObserver Creates YArrayEvent and calls observers.
func (b YArrayBase) CallObserver(transaction *utils.Transaction, parentSubs map[string]struct{}) {
	if !transaction.Local {
		b.searchMarkers.Clear()
	}
}

func (b YArrayBase) InsertGenericsAfter(transaction *utils.Transaction, referenceItem *structs.Item, contents []any) {
	var left = referenceItem
	var doc = transaction.Doc
	var ownClientId = doc.ClientId
	var store = doc.Store
	var right *structs.Item
	if referenceItem == nil {
		right = b.Start
	} else {
		right = referenceItem.Right.(*structs.Item)
	}
	var jsonContent []any

	packJsonContent := func() {
		if len(jsonContent) > 0 {
			left = &structs.Item{
				Id:          &utils.ID{Client: ownClientId, Clock: store.GetState(ownClientId)},
				Left:        left,
				LeftOrigin:  left.LastId,
				Right:       right,
				RightOrigin: right.Id,
				Parent:      b,
				ParentSub:   "",
				Content:     content.NewAny(jsonContent),
			}
			left.Integrate(transaction, 0)
			jsonContent = []any{}
		}
	}
	for _, c := range contents {
		left = &structs.Item{
			Id:          &utils.ID{Client: ownClientId, Clock: store.GetState(ownClientId)},
			Left:        left,
			LeftOrigin:  left.LastId,
			Right:       right,
			RightOrigin: right.Id,
			Parent:      b,
			ParentSub:   "",
		}
		switch a := c.(type) {
		case []byte:
			packJsonContent()
			left.Content = content.NewBinary(a)
			left.Integrate(transaction, 0)
			break
		case *utils.YDoc:
			packJsonContent()
			left.Content = content.NewDoc(a)
			left.Integrate(transaction, 0)
			break
		case AbstractType:
			packJsonContent()
			left.Content = content.NewType(a)
			left.Integrate(transaction, 0)
			break
		default:
			jsonContent = append(jsonContent, c)
			break
		}
	}

	packJsonContent()
}

// FindMarker Search markers help us to find positions in the associative array faster.
// They speed up the process of finding a position without much bookkeeping.
// A maximum of 'MaxSearchMarker' objects are created.
// This function always returns a refreshed marker (updated timestamp).
func (b YArrayBase) FindMarker(index uint64) *ArraySearchMarker {
	if b.Start == nil || index == 0 || b.searchMarkers == nil || b.searchMarkers.Count() == 0 {
		return nil
	}

	var Aggregate = func(ac *ArraySearchMarkerCollection) *ArraySearchMarker {
		var ret = ac.searchMarkers[0]
		for i, marker := range ac.searchMarkers {
			if i == 0 {
				continue
			}
			if math.Abs(float64(index-marker.Index)) < math.Abs(float64(index-ret.Index)) {
				ret = marker
			}
		}
		return ret
	}
	var marker *ArraySearchMarker
	if b.searchMarkers.Count() != 0 {
		marker = Aggregate(b.searchMarkers)
	}
	var p = b.Start

	var pIndex uint64

	if marker != nil {
		p = marker.P
		pIndex = marker.Index

		// We used it, we might need to use it again.
		marker.RefreshTimestamp()
	}

	// Iterate to right if possible.
	for p.Right != nil && pIndex < index {
		if !p.Deleted && p.Countable {
			if index < pIndex+p.Length {
				break
			}
			pIndex += p.Length
		}
		p = p.Right.(*structs.Item)
	}

	// Iterate to left if necessary (might be that pIndex > index).
	for p.Left != nil && pIndex > index {
		p = p.Left.(*structs.Item)
		if p == nil {
			break
		} else if !p.Deleted && p.Countable {
			pIndex -= p.Length
		}
	}

	// We want to make sure that p can't be merged with left, because that would screw up everything.
	// In that case just return what we have (it is most likely the best marker anyway).
	// Iterate to left until p can't be merged with left.
	for p.Left != nil && p.Left.(structs.IAbstractStruct).ID().Client == p.Id.Client &&
		p.Left.(structs.IAbstractStruct).ID().Clock+p.Left.(structs.IAbstractStruct).GetLength() == p.Id.Clock {
		p = p.Left.(*structs.Item)
		if p == nil {
			break
		} else if !p.Deleted && p.Countable {
			pIndex -= p.Length
		}
	}

	if marker != nil && uint64(math.Abs(float64(marker.Index)-float64(pIndex))) <
		p.Parent.(AbstractType).Length/uint64(MaxSearchMarkers) {
		// Adjust existing marker.
		marker.Update(p, pIndex)
		return marker
	}
	// Create a new marker.
	return b.searchMarkers.MarkPosition(p, pIndex)

}

func (b YArrayBase) InternalSlice(start, end uint64) []any {
	if start < 0 {
		start += b.Length
	}
	if end < 0 {
		end += b.Length
	}

	if start < 0 {
		return nil
		// throw
		// new
		// ArgumentOutOfRangeException(nameof(start))
	}
	if end < 0 {
		return nil
	}
	if start > end {
		return nil
	}

	length := end - start
	// Debug.Assert(length >= 0)

	var (
		cs []any
		n  = b.Start
	)
	for n != nil && length > 0 {
		if n.Countable && !n.Deleted {
			var c = n.Content.GetContent().([]any)
			if uint64(len(c)) <= start {
				start -= uint64(len(c))
			} else {
				for i := start; i < uint64(len(c)) && length > 0; i++ {
					cs = append(cs, c[i])
					length--
				}

				start = 0
			}
		}

		n = n.Right.(*structs.Item)
	}

	return cs
}
