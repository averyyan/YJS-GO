package encoder

import (
	"bufio"

	"YJS-GO/lib0"
)

var _ IEncoder[uint] = (*IncUintOptRleEncoder)(nil)

type IncUintOptRleEncoder struct {
	AbstractEncoder

	state uint
	count uint

	writer *bufio.Writer
}

func (e *IncUintOptRleEncoder) Write(value any) {
	e.CheckDisposed()

	if e.state+e.count == value {
		e.count++
	} else {
		e.WriteEncodedValue()

		e.count = 1
		e.state = value.(uint)
	}
}

func (e *IncUintOptRleEncoder) WriteEncodedValue() {
	if e.count > 0 {
		// Flush counter, unless this is the first value (count = 0).
		// Case 1: Just a single value. Set sign to positive.
		// Case 2: Write several values. Set sign to negative to indicate that there is a length coming.
		if e.count == 1 {
			lib0.WriteVarInt2(e.writer, int(e.state))
		} else {
			// Specify 'treatZeroAsNegative' in case we pass the '-0' value.
			lib0.WriteVarInt(e.writer, int(-e.state), e.state == 0)

			// Since count is always >1, we can decrement by one. Non-standard encoding.
			lib0.WriteVarUint(e.writer, e.count-2)
		}
	}
}
