package decoder

import (
	"bufio"
	"encoding/binary"
)

type UintOptRleDecoder struct {
	Reader *bufio.Reader
	State  uint64
	Count  uint64
}

func (d UintOptRleDecoder) Read() uint64 {
	// d.CheckDisposed()

	if d.Count == 0 {
		binary.ReadVarint(d.Reader)
		var (
			value, sign
		) = Stream.ReadVarInt()

		// If the sign is negative, we read the count too; otherwise, count is 1.
		bool
		isNegative = sign < 0
		if isNegative {
			_state = (uint)(-value)
			_count = Stream.ReadVarUint() + 2
		} else
		{
			_state = (uint)
			value
			_count = 1
		}
	}

	_count--
	return _state
}
