package decoder

import (
	"bufio"
	"encoding/binary"
	"errors"
	"io"

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
		value, sign, err := ReadVarInt(i.reader)
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

func ReadVarInt(reader io.Reader) (uint, uint, error) {
	byteReader := bufio.NewReader(reader)
	r, err := byteReader.ReadByte()
	if err != nil {
		return 0, 0, err
	}
	var num uint = uint(r) & lib0.Bits6
	var len = 6
	var sign uint = 1
	if uint(r)&lib0.Bit7 > 0 {
		sign = -1
	}

	if (uint(r) & lib0.Bit8) == 0 {
		// Don't continue reading.
		return sign * num, sign, nil
	}
	for true {
		r, err := byteReader.ReadByte()
		if err != nil {
			return 0, 0, err
		}
		num |= (uint(r) & lib0.Bits7) << len
		len += 7

		if uint(r) < lib0.Bit8 {
			return sign * num, sign, nil
		}

		if len > 41 {
			// throw new InvalidDataException("Integer out of range")
			return 0, 0, errors.New("Integer out of range")
		}
	}
	return 0, 0, err
}
