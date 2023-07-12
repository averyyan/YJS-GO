package utils

import (
	"container/list"

	"YJS-GO/structs"
	"YJS-GO/types"
)

const (
	Add = iota + 1
	Update
	Delete
)

type ChangeAction int

type ChangesCollection struct {
	Added   map[structs.Item]struct{}
	Deleted map[structs.Item]struct{}
	Delta   list.List
	Keys    map[string]ChangeKey
}

type Delta struct {
	Insert     interface{}
	Delete     int
	Retain     int
	Attributes map[string]interface{}
}

type ChangeKey struct {
	Action   ChangeAction
	OldValue interface{}
}

type NewBaseType struct {
	changes       ChangesCollection
	Target        *types.AbstractType
	CurrentTarget *types.AbstractType
	Transaction   Transaction
	Path          []any
}

type YEvent struct {
	NewBaseType
}
