package utils

import (
	"bufio"
)

type IUpdateDecoder interface {
	IDSDecoder
	ReadLeftId() *ID
	ReadRightId() *ID
	ReadClient() uint64
	ReadInfo() byte
	ReadString() string
	ReadParentInfo() bool
	ReadTypeRef() uint
	ReadLength() uint64
	ReadAny() any
	ReadBuffer() []byte
	ReadKey() string
	ReadJson() any
}

type IDSDecoder interface {
	Reader() *bufio.Reader
	ReadDsLength() uint64
	ReadDsClock() uint64
	ResetDsCurVal()
}
