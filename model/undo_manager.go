package model

import (
	"container/list"
	"time"

	"YJS-GO/structs"
	"YJS-GO/types"
	"YJS-GO/utils"
)

type OperationType int

const (
	Undo = iota + 1
	Redo
)

type StackItem struct {
	BeforeState map[int]int
	AfterState  map[int]int
	Meta        map[string]any
	DeleteSet   utils.DeleteSet
}

type StackEventArgs struct {
	StackItem          StackItem
	Type               OperationType
	ChangedParentTypes map[*types.AbstractType][]YEvent // public IDictionary<AbstractType,
	// IList<YEvent>> ChangedParentTypes { get; }
	Origin any
}

type UndoManager struct {
	Scope          []any
	DeleteFilter   func(item structs.Item) bool
	TrackedOrigins map[any]struct{}
	UndoStack      list.List
	RedoStack      list.List
	// Whether the client is currently undoing (calling UndoManager.Undo()).
	Undoing        bool
	Redoing        bool
	Doc            YDoc
	LastChange     time.Time
	CaptureTimeout int
}

func (*UndoManager) undo() {

}
func (*UndoManager) redo() {

}
func (*UndoManager) stopCapturing() {

}
