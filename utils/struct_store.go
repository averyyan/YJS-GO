package utils

import (
	"container/list"

	"YJS-GO/structs"
)

var Clients = map[uint64]list.List{}
var _pendingClientStructRefs = map[uint64]PendingClientStructRef{}
var _pendingStack = list.New()

type StructStore struct {
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
