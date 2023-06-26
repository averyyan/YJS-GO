package model

import "YJS-GO/structs"

type transaction[K string, V any] struct {
	beforeState        map[int]int
	changedParentTypes map[ytype][]YEvent
	changedTypes       map[ytype]string
	deletedStructs     map[structs.item]struct{}
	newTypes           map[structs.item]struct{}
	y                  Y[K, V]
}
