package types

import (
	"math"

	"YJS-GO/structs"
	"github.com/uber-go/atomic"
)

type YArray struct {
}

// Assigned to '-1', so the first timestamp is '0'.
var globalSearchMarkerTimestamp = atomic.NewInt64(-1)

type ArraySearchMarker struct {
	P         *structs.Item
	Index     uint64
	Timestamp int64
}

func NewArraySearchMarker(p *structs.Item, index uint64) *ArraySearchMarker {
	t := &ArraySearchMarker{
		P:     p,
		Index: index,
	}
	t.P.Marker = true
	t.RefreshTimestamp()
	return t
}

func (a ArraySearchMarker) RefreshTimestamp() {
	a.Timestamp = globalSearchMarkerTimestamp.Inc()
}

func (a ArraySearchMarker) Update(p *structs.Item, index uint64) {
	a.P.Marker = false
	a.P = p
	a.P.Marker = true
	a.Index = index
	a.RefreshTimestamp()
}

var MaxSearchMarkers = 80

type ArraySearchMarkerCollection struct {
	searchMarkers []*ArraySearchMarker
	count         uint64
}

func (ac ArraySearchMarkerCollection) Clear() {
	ac.searchMarkers = []*ArraySearchMarker{}
}

func (ac ArraySearchMarkerCollection) MarkPosition(p *structs.Item, index uint64) *ArraySearchMarker {
	if len(ac.searchMarkers) >= MaxSearchMarkers {
		// Override oldest marker (we don't want to create more objects).
		var marker = minSearchMarks(ac.searchMarkers)
		marker.Update(p, index)
		return marker
	} else {
		// Create a new marker.
		var pm = NewArraySearchMarker(p, index)
		ac.searchMarkers = append(ac.searchMarkers, pm)
		return pm
	}
}

// UpdateMarkerChanges Update markers when a change happened.
// This should be called before doing a deletion!
func (ac ArraySearchMarkerCollection) UpdateMarkerChanges(index, length uint64) {
	for i := len(ac.searchMarkers) - 1; i >= 0; i-- {
		var m = ac.searchMarkers[i]

		if length > 0 {
			var p = m.P
			p.Marker = false
			// Ideally we just want to do a simple position comparison, but this will only work if
			// search markers don't point to deleted items for formats.
			// Iterate marker to prev undeleted countable position so we know what to do when updating a position.
			for p != nil && (p.Deleted || !p.Countable) {
				// Debug.Assert(p.Left != p)
				p = p.Left.(*structs.Item)
				if p != nil && !p.Deleted && p.Countable {
					// Adjust position. The loop should break now.
					m.Index -= p.Length
				}
			}

			if p == nil || p.Marker {
				// Remove search marker if updated position is nil or if position is already marked.
				ac.searchMarkers = append(ac.searchMarkers[:i], ac.searchMarkers[i+1:]...)
				continue
			}

			m.P = p
			p.Marker = true
		}

		// A simple index <= m.Index check would actually suffice.
		if index < m.Index || (length > 0 && index == m.Index) {
			m.Index = uint64(math.Max(float64(index), float64(m.Index+length)))
		}
	}
}

func (ac ArraySearchMarkerCollection) Count() int {
	return len(ac.searchMarkers)
}

func minSearchMarks(asm []*ArraySearchMarker) *ArraySearchMarker {
	if len(asm) == 0 {
		return nil
	}
	if len(asm) == 1 {
		return asm[0]
	}
	min := asm[0]
	for i, marker := range asm {
		if i == 0 {
			continue
		}
		if min.Timestamp < marker.Timestamp {
			min = marker
		}
	}
	return min
}
