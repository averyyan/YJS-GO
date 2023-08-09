package decoder

import (
	"bufio"
	"encoding/binary"
)

type IntDiffOptRleDecoder struct {
	Reader *bufio.Reader
	State  uint64
}

func (d IntDiffOptRleDecoder) Read() uint64 {
	varint, err := binary.ReadVarint(d.Reader)
	if err != nil {
		return 0
	}
	d.State = uint64(varint)
	return d.State
}
