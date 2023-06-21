package utils

import (
	"encoding/binary"
	"io"
)

type ID struct {
	client uint64
	clock  uint64
}

func (a *ID) EQ(b *ID) bool {
	return a == b || (a.client == b.client && a.clock == b.clock)
}

func (a *ID) writeID(buffer []byte) {
	binary.AppendUvarint(buffer, a.client)
	binary.AppendUvarint(buffer, a.clock)
}

func ReadID(reader io.ByteReader) *ID {
	client, err := binary.ReadUvarint(reader)
	if err != nil {
		return nil
	}
	clock, err := binary.ReadUvarint(reader)
	if err != nil {
		return nil
	}
	return &ID{
		client: client,
		clock:  clock,
	}
}
