package encoder

import (
	"bufio"

	"YJS-GO/lib0"
)

var _ IEncoder[uint64] = (*UintOptRleEncoder)(nil)

// UintOptRleEncoder Optimized RLE encoder that does not suffer from the mentioned problem of the basic RLE encoder.
// Internally uses VarInt encoder to write unsigned integers.
// If the input occurs multiple times, we write it as a negative number. The <see cref="UintOptRleDecoder"/>
// then understands that it needs to read a count.
type UintOptRleEncoder struct {
	AbstractEncoder
	State  uint64
	Count  uint64
	Writer *bufio.Writer
}

func (e *UintOptRleEncoder) Write(v any) {
	e.CheckDisposed()
	if e.State == v {
		e.Count++
	} else {
		e.WriteEncodedValue()
		e.Count = 1
		e.State = v.(uint64)
	}
}

func (e *UintOptRleEncoder) Flush() {
	e.WriteEncodedValue()
	e.FlushV()
}

func (e *UintOptRleEncoder) WriteEncodedValue() {
	if e.Count > 0 {
		// Flush counter, unless this is the first value (count = 0).
		// Case 1: Just a single value. Set sign to positive.
		// Case 2: Write several values. Set sign to negative to indicate that there is a length coming.
		if e.Count == 1 {
			lib0.WriteVarInt2(e.Writer, int(e.State))
		} else {
			// Specify 'treatZeroAsNegative' in case we pass the '-0'.
			lib0.WriteVarInt(e.Writer, int(-e.State), e.State == 0)
			// Since count is always >1, we can decrement by one. Non-standard encoding.
			lib0.WriteVarUint(e.Writer, uint(e.Count-2))
		}
	}
}
