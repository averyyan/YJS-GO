package utils

import (
	"container/list"

	"YJS-GO/structs"
)

var Clients = map[uint64]list.List{}
var _pendingClientStructRefs = map[uint64]PendingClientStructRef{}
var _pendingStack = list.New()

type StructStore struct {
	Clients map[uint64][]structs.IAbstractStruct
}

func FindIndexSS(abstractStructs []structs.IAbstractStruct, clock uint64) uint {
	if len(abstractStructs) <= 0 {
		return 0
	}
	var left = 0
	var right = len(abstractStructs) - 1
	var mid = abstractStructs[right].(*structs.Item)
	var midClock = mid.Id.Clock

	if midClock == clock {
		return uint(right)
	}

	// @todo does it even make sense to pivot the search?
	// If a good split misses, it might actually increase the time to find the correct item.
	// Currently, the only advantage is that search with pivoting might find the item on the first try.
	midIndex := (int)((clock * uint64(right)) / (midClock + mid.Length - 1))

	for left <= right {
		mid = abstractStructs[midIndex].(*structs.Item)
		midClock = mid.Id.Clock

		if midClock <= clock {
			if clock < midClock+mid.Length {
				return uint(midIndex)
			}

			left = midIndex + 1
		} else {
			right = midIndex - 1
		}

		midIndex = (left + right) / 2
	}
	// Always check state before looking for a struct in StructStore
	// Therefore the case of not finding a struct is unexpected.
	return 0
}

type PendingClientStructRef struct {
	NextReadOperation int
	Ref               list.List
}

// / <summary>
// / Return the states as a Map<int,int>.
// / Note that Clock refers to the next expected Clock id.
// / </summary>
func (s StructStore) GetStateVector() map[uint64]uint64 {
	var result = map[uint64]uint64{}

	// foreach (var kvp in Clients)
	//            {
	//                var str = kvp.Value[kvp.Value.Count - 1];
	//                result[kvp.Key] = str.Id.Clock + str.Length;
	//            }
	for k, v := range Clients {
		var str = v.Back().Value.(*structs.AbstractStruct)
		result[k] = str.Id.Clock + uint64(str.Length)
	}
	return result
}

func (s StructStore) MergeReadStructsIntoPendingReads(refs map[int]list.List) {
	// TODO
}

func (s StructStore) ResumeStructIntegration(transaction Transaction) {
	// todo
}

func (s StructStore) CleanupPendingStructs() {
	// todo
}

func (s StructStore) TryResumePendingDeleteReaders(transaction Transaction) {
	// todo
}

func (s StructStore) GetState(c uint64) uint64 {
	ss, ok := s.Clients[c]
	if ok {
		lastStruct := ss[len(ss)-1].(*structs.Item)
		return lastStruct.Id.Clock + lastStruct.Length
	}
	return 0
}
