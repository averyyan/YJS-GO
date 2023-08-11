package decoder

import (
	"bufio"

	"YJS-GO/lib0"
)

var _ IDecoder[byte] = (*RleDecoder)(nil)

type RleDecoder struct {
	AbstractDecoder
	Reader *bufio.Reader
	State  byte
	Count  int
}

func (r RleDecoder) Read() byte {
	r.CheckDisposed()
	var err error
	if r.Count == 0 {
		r.State, err = r.Reader.ReadByte()
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
