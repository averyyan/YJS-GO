package utils

import (
	"bufio"
	"bytes"
	"encoding/binary"

	"YJS-GO/lib0"
)

type RelativePosition struct {
	ItemId *ID
	TypeId *ID
	TName  string
	// A relative position is associated to a specific character.
	// By default, the value is <c>&gt;&eq; 0</c>, the relative position is associated to the character
	// after the meant position.
	// I.e. position <c>1</c> in <c>'ab'</c> is associated with the character <c>'b'</c>.
	// If the value is <c>&lt; 0</c>, then the relative position is associated with the character
	// before the meant position.
	Assoc int
}

func ReadByteRP(b []byte) *RelativePosition {
	return ReadRP(bufio.NewReader(bytes.NewReader(b)))
}

func ReadRP(reader *bufio.Reader) *RelativePosition {
	a, err := binary.ReadUvarint(reader)
	if err != nil {
		return nil
	}
	var (
		itemID *ID
		tname  string
		typeID *ID
	)
	switch a {
	case 0:
		// Case 1: Found position somewhere in the linked list.
		itemID = ReadID(reader)
	case 2:
		// Case 2: Found position at the end of the list and type is stored in y.share.
		typeID = ReadID(reader)
	case 1:
		// Case 3: Found position at the end of the list and type is attached to an item.
		tname = lib0.ReadVarString(reader)
	default:
		return nil
	}
	rp := &RelativePosition{
		ItemId: itemID,
		TypeId: typeID,
		TName:  tname,
	}
	if reader.Buffered() > 0 {
		assoc, err := binary.ReadVarint(reader)
		if err != nil {
			return nil
		}
		rp.Assoc = int(assoc)
	}
	return rp
}
