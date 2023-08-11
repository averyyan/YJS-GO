package decoder

import (
	"bufio"

	"YJS-GO/lib0"
)

var _ IDecoder[int] = (*IntDiffOptRleDecoder)(nil)

type IntDiffOptRleDecoder struct {
	AbstractDecoder
	Reader *bufio.Reader
	State  int
	Count  uint64
	Diff   int
}

func (d IntDiffOptRleDecoder) Read() int {
	d.CheckDisposed()

	if d.Count == 0 {
		_, value, _ := lib0.ReadVarInt(d.Reader)
		var diff = value

		// If the first bit is set, we read more data.
		hasCount := (diff & lib0.Bit1) > 0

		if diff < 0 {
			d.Diff = -((-int(diff)) >> 1)
		} else {
			d.Diff = int(diff) >> 1
		}
		if hasCount {
			d.Count = uint64(lib0.ReadVarUint(d.Reader) + 2)
		} else {
			d.Count = 1
		}
	}

	d.State += d.Diff
	d.Count--
	return d.State
}
