package utils

import "io"

type IUpdateDecoder interface {
	Reader() io.ByteReader
	// parent interface
	ReadDsLength()
	ReadDsClock()
	ResetDsCurVal()
	// self interface
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
