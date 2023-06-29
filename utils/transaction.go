package utils

import "YJS-GO/structs"

type Transaction struct {
	Doc            YDoc
	BeforeState    map[int]int
	DeletedStructs map[structs.Item]struct{}
	NewTypes       map[structs.Item]struct{}
}
