package lib0

import (
	"bufio"
	"bytes"
	"encoding/binary"
	"io"
)

type VarReader struct {
	*bufio.Reader
}

func New(reader io.Reader) *VarReader {
	return &VarReader{Reader: bufio.NewReader(reader)}
}

func ReadVarString(reader *bufio.Reader) string {
	remainingLen, err := binary.ReadUvarint(reader)
	if err != nil || remainingLen == 0 {
		return ""
	}
	readBytes := make([]byte, remainingLen)
	n, err := reader.Read(readBytes)
	if err != nil {
		return ""
	}
	if uint64(n) != remainingLen {
		return ""
	}

	// reader, _ := charset.NewReaderLabel(origEncoding, byteReader)
	// strBytes, _ = ioutil.ReadAll(reader)
	// return string(strBytes)
	return string(bytes.Runes(readBytes))
}
