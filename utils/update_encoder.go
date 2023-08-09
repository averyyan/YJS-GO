package utils

import (
	"bufio"
)

type IDSEncoder interface {
	RestWriter() *bufio.ReadWriter

	ToArray() []byte
	WriteDsLength(length uint64)
	WriteDsClock(clock uint64)
	ResetDsCurVal()
}

type IUpdateEncoder interface {
	IDSEncoder
	Writer() *bufio.Writer

	WriteLeftId(id *ID)
	WriteRightId(id *ID)

	// WriteClient NOTE: Use 'writeClient' and 'writeClock' instead of writeID if possible.
	WriteClient(client uint64)
	WriteInfo(info byte)
	WriteString(s string)
	WriteParentInfo(isYKey bool)
	WriteTypeRef(info uint)

	// WriteLength Write len of a struct - well suited for Opt RLE encoder.
	WriteLength(len int)
	WriteAny(object any)
	WriteBuffer(buf []byte)
	WriteKey(key string)
	WriteJson(T any)
}
