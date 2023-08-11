package decoder

import (
	"bufio"

	"YJS-GO/lib0"
)

var _ IDecoder[int] = (*RleIntDiffDecoder)(nil)

type RleIntDiffDecoder struct {
	AbstractDecoder
	Reader *bufio.Reader
	State  int
	Count  int
}

func (r RleIntDiffDecoder) Read() int {
	r.CheckDisposed()
	var err error
	if r.Count == 0 {
		_, value, _ := lib0.ReadVarInt(r.Reader)
		r.State += int(value)
		if err != nil {
			return 0
		}
		if r.HasContent() {
			// See encoder implementation for the reason why this is incremented.
			r.Count = int(lib0.ReadVarUint(r.Reader) + 1)
			// Debug.Assert(r.count > 0)
		} else {
			// Read the current value forever.
			r.Count = -1
		}
	}

	r.Count--
	return r.State
}
