package encoder

import (
	"bufio"

	"YJS-GO/lib0"
)

var _ IEncoder[uint64] = (*IntDiffOptRleEncoder)(nil)

// IntDiffOptRleEncoder A combination of the <see cref="IntDiffEncoder"/> and the <see cref="UintOptRleEncoder"/>.
// The count approach is similar to the <see cref="UintOptRleDecoder"/>, but instead of using
// the negative bitflag, it encodes in the LSB whether a count is to be read.
// <br/>
// WARNING: Therefore this encoder only supports 31 bit integers.
// <br/>
// Encodes <c>[1, 2, 3, 2]</c> as <c>[3, 1, -2]</c> (more specifically <c>[(1 << 1) | 1, (3 << 0) | 0, -((1 << 1) | 0)]</c>).
// <br/>
// Internally uses variable length encoding. Contrary to the normal UintVar encoding, the first byte contains:
// * 1 bit that denotes whether the next value is a count (LSB).
// * 1 bit that denotes whether this value is negative (MSB - 1).
// * 1 bit that denotes whether to continue reading the variable length integer (MSB).
// <br/>
// Therefore, only five bits remain to encode diff ranges.
// <br/>
// Use this encoder only when appropriate. In most cases, this is probably a bad idea.
// </summary>
// <seealso cref="IntDiffOptRleDecoder"/>
type IntDiffOptRleEncoder struct {
	AbstractEncoder
	state  uint64
	diff   uint64
	count  uint
	writer *bufio.Writer
}

func (e *IntDiffOptRleEncoder) Write(value any) {
	// Debug.Assert(value <= lib0.Bits30)
	e.CheckDisposed()

	if e.diff == value.(uint64)-e.state {
		e.state = value.(uint64)
		e.count++
	} else {
		e.WriteEncodedValue()

		e.count = 1
		e.diff = value.(uint64) - e.state
		e.state = value.(uint64)
	}
}

func (e *IntDiffOptRleEncoder) Flush() {
	e.WriteEncodedValue()
	e.FlushV()
}

func (e *IntDiffOptRleEncoder) WriteEncodedValue() {
	if e.count > 0 {
		var encodedDiff uint64
		if e.diff < 0 {
			tmp := 1
			if e.count == 1 {
				tmp = 0
			}
			encodedDiff = uint64(-(((uint)(-e.diff) << 1) | (uint)(tmp)))
		} else {
			// 31bit making up a diff  | whether to write the counter.
			tmp := 1
			if e.count == 1 {
				tmp = 0
			}
			encodedDiff = e.diff<<1 | (uint64)(tmp)
		}
		lib0.WriteVarInt2(e.writer, int(encodedDiff))
		if e.count > 1 {
			// Since count is always >1, we can decrement by one. Non-standard encoding.
			lib0.WriteVarUint(e.writer, e.count-2)
		}
	}
}
