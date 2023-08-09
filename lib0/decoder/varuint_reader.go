package decoder

import (
	"bufio"
	"bytes"
	"encoding/binary"
)

func ReadVarUint8ArrayAsStream(reader *bufio.Reader) *bufio.Reader {
	var data = ReadVarUint8Array(reader)
	return bufio.NewReader(bytes.NewReader(data))
}

func ReadVarUint8Array(reader *bufio.Reader) []byte {
	// uint len = stream.ReadVarUint();
	length, err := binary.ReadUvarint(reader)
	if err != nil {
		return nil
	}
	var ret = make([]byte, length)
	n, err := reader.Read(ret)
	if err != nil || uint64(n) != length {
		return nil
	}
	return ret
}
