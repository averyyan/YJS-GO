package utils

import "io"

type IUpdateDecoder interface {
	IDSDecoder
	ReadLeftId() ID
	ReadRightId() ID
	ReadClient() int64
	ReadInfo() byte
	ReadString() string
	ReadParentInfo() bool
	ReadTypeRef() uint
	ReadLength() int
	ReadAny() any
	ReadBuffer() []byte
	ReadKey() string
	ReadJson() any
}

type IDSDecoder interface {
	Reader() io.Reader
	ReadDsLength()
	ReadDsClock()
	ResetDsCurVal()
}
