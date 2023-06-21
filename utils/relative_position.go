package utils

import (
	"bufio"
	"bytes"
	"encoding/binary"

	"YJS-GO/lib0"
)

type RelativePosition struct {
	Item   *ID
	TypeId *ID
	TName  string
	Assoc  int
}

func ReadByteRP(b []byte) *RelativePosition {
	return ReadRP(bufio.NewReader(bytes.NewReader(b)))
}

func ReadRP(reader *bufio.Reader) *RelativePosition {
	a, err := binary.ReadUvarint(reader)
	if err != nil {
		return nil
	}
	var itemID *ID
	var tname string
	var typeID *ID
	switch a {
	case 0:
		typeID = ReadID(reader)
	case 2:
		itemID = ReadID(reader)
	case 1:
		tname = lib0.ReadVarString(reader)
	}
	return &RelativePosition{
		Item:   itemID,
		TypeId: typeID,
		TName:  tname,
		Assoc:  0,
	}
}
