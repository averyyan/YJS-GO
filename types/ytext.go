package types

import "YJS-GO/utils"

const (
	Insert = iota
	Delete
	Retain
)

type YText struct {
	YArrayBase
}

type YTextEvent struct {
	utils.YEvent
	subs  map[string]struct{}
	delta []*utils.Delta
}

func NewYTextEvent(arr Ytext) {

}
