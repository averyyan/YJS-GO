package decoder

import (
	"bufio"
	"encoding/binary"

	"YJS-GO/lib0"
)

var _ IDecoder[uint] = (*IncUintOptRleDecoder)(nil)

type IncUintOptRleDecoder struct {
	state     uint64
	count     uint64
	leaveOpen bool
	reader    *bufio.Reader
	Disposed  bool
}

func (i IncUintOptRleDecoder) Read(p []byte) (n int, err error) {
	if err = i.CheckDisposed(); err != nil {

	}
	if i.count == 0 {
		value, sign, err := lib0.ReadVarInt(i.reader)
		if err != nil {

		}

		// If the sign is negative, we read the count too; otherwise. count is 1.
		isNegative := sign < 0
		if isNegative {
			i.state = uint64(-value)
			tmp, _ := binary.ReadUvarint(i.reader)
			i.count = tmp + 2
		} else {
			i.state = uint64(value)
			i.count = 1
		}
	}
	i.count--
	return int(i.state), nil
}

func (i IncUintOptRleDecoder) CheckDisposed() error {
	return nil
}

func (i IncUintOptRleDecoder) ReadV() any {
	// TODO implement me
	panic("implement me")
}

func NewIncUintOptRleDecoder(input *bufio.Reader, leaveOpen bool) *IncUintOptRleDecoder {
	return &IncUintOptRleDecoder{
		leaveOpen: leaveOpen,
		reader:    input,
	}
}
