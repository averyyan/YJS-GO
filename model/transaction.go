package model

type transaction[K string, V any] struct {
	beforeState        map[int]int
	changedParentTypes map[ytype][]YEvent
	changedTypes       map[ytype]string
	deletedStructs     map[item]struct{}
	newTypes           map[item]struct{}
	y                  Y[K, V]
}
