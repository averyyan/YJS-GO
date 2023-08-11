package decoder

import (
	"bufio"
	"encoding/binary"

	"YJS-GO/lib0"
)

var _ IDecoder[uint64] = (*UintOptRleDecoder)(nil)

type UintOptRleDecoder struct {
	AbstractDecoder
	Reader *bufio.Reader
	State  uint64
	Count  uint64
}

func (d UintOptRleDecoder) Read() uint64 {
	d.CheckDisposed()

	if d.Count == 0 {
		binary.ReadVarint(d.Reader)
		var value, sign, err = lib0.ReadVarInt(d.Reader)
		if err != nil {
			return 0
		}
		// If the sign is negative, we read the count too; otherwise, count is 1.
		isNegative := sign < 0
		if isNegative {
			d.State = uint64(-value)
			d.Count = uint64(lib0.ReadVarUint(d.Reader) + 2)
		} else {
			d.State = uint64(value)
			d.Count = 1
		}
	}

	d.Count--
	return d.State
}
