package utils

import "YJS-GO/structs"

type AbsolutePosition struct {
	Type  structs.AbstractStruct
	Index int
	Assoc int
	store StructStore
}

func (ap AbsolutePosition) TryCreateFromRelativePosition() {

}
