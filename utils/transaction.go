package utils

import (
	"YJS-GO/structs"
	"YJS-GO/types"
)

type Transaction struct {
	Doc            YDoc
	BeforeState    map[uint64]uint64
	DeletedStructs map[*structs.Item]struct{}
	NewTypes       map[*structs.Item]struct{}
	DeleteSet      DeleteSet
	Changed        map[*types.AbstractType]map[string]struct{}
}

func (t Transaction) AddChangedTypeToTransaction(ty *types.AbstractType, parentSub string) {
	var item = ty.Item

	clock, ok := t.BeforeState[item.Id.Client]
	if item == nil || (ok && item.Id.Clock < clock && !item.Deleted) {
		var set map[string]struct{}
		set, ok := t.Changed[ty]
		if !ok {
			set = map[string]struct{}{}
			t.Changed[ty] = set
		}
		set[parentSub] = struct{}{}
	}
}
