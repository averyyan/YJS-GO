package encoder

import (
	"bufio"

	"YJS-GO/lib0"
)

var _ IEncoder[byte] = (*RleEncoder)(nil)

type RleEncoder struct {
	AbstractEncoder
	writer *bufio.Writer
	state  byte
	count  uint64
}

func (r *RleEncoder) Write(value any) {
	r.CheckDisposed()

	if r.state == value {
		r.count++
	} else {
		if r.count > 0 {
			// Flush counter, unless this is the first value (count = 0).
			// Since 'count' is always >0, we can decrement by one. Non-standard encoding.
			lib0.WriteVarUint(r.writer, uint(r.count-1))
		}

		err := r.writer.WriteByte(value.(byte))
		if err != nil {
			return
		}

		r.count = 1
		r.state = value.(byte)
	}
}
