package utils

import "container/list"

var Clients = map[uint64]list.List{}
var _pendingClientStructRefs = map[uint64]PendingClientStructRef{}
var _pendingStack = list.New()

type StructStore struct {
}

type PendingClientStructRef struct {
	NextReadOperation int
	Ref               list.List
}
