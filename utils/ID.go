package utils

import (
	"encoding/binary"
	"io"
)

type ID struct {
	Client uint64
	Clock  uint64
}

func (a *ID) EQ(b *ID) bool {
	return a == b || (a.Client == b.Client && a.Clock == b.Clock)
}

func (a *ID) writeID(buffer []byte) {
	binary.AppendUvarint(buffer, a.Client)
	binary.AppendUvarint(buffer, a.Clock)
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
		Client: client,
		Clock:  clock,
	}
}
