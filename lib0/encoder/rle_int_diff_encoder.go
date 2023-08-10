package encoder

import (
	"bufio"

	"YJS-GO/lib0"
)

var _ IEncoder[uint64] = (*RleIntDiffEncoder)(nil)

// RleIntDiffEncoder A combination of <see cref="IntDiffEncoder"/> and <see cref="RleEncoder"/>.
// Basically first writes the <see cref="IntDiffEncoder"/> and then counts duplicate
// diffs using the <see cref="RleEncoder"/>.
// Encodes values <c>[1, 1, 1, 2, 3, 4, 5, 6]</c> as <c>[1, 1, 0, 2, 1, 5]</c>.
type RleIntDiffEncoder struct {
	AbstractEncoder
	writer *bufio.Writer
	state  uint64
	count  uint
}

func NewRleIntDiffEncoder(start uint64) *RleIntDiffEncoder {
	a := &RleIntDiffEncoder{
		AbstractEncoder: AbstractEncoder{},
		state:           start,
	}
	return a
}

func (r *RleIntDiffEncoder) Write(value any) {
	r.CheckDisposed()

	if r.state == value && r.count > 0 {
		r.count++
	} else {
		if r.count > 0 {
			lib0.WriteVarUint(r.writer, r.count-1)
		}

		lib0.WriteVarInt2(r.writer, int(value.(uint64)-r.state))

		r.count = 1
		r.state = value.(uint64)
	}
}
