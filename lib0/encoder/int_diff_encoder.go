package encoder

import (
	"bufio"

	"YJS-GO/lib0"
)

var _ IEncoder[uint64] = (*IntDiffEncoder)(nil)

type IntDiffEncoder struct {
	AbstractEncoder
	Writer *bufio.Writer
	State  uint64
}

func NewIntDiffEncoder(start uint64) *IntDiffEncoder {
	t := &IntDiffEncoder{
		AbstractEncoder: AbstractEncoder{},
		Writer:          nil,
		State:           start,
	}
	return t
}

// <inheritdoc/>
func (e *IntDiffEncoder) Write(value any) {
	e.CheckDisposed()

	lib0.WriteVarInt2(e.Writer, int(value.(uint64)-e.State))
	e.State = value.(uint64)
}
