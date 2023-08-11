package decoder

import (
	"bufio"

	"YJS-GO/lib0"
)

var _ IDecoder[int] = (*IntDiffDecoder)(nil)

// IntDiffDecoder Basic diff encoder using variable length encoding.
// Encodes the values <c>[3, 1100, 1101, 1050, 0]</c> to <c>[3, 1097, 1, -51, -1050]</c>.
// <seealso cref="IntDiffDecoder"/>
type IntDiffDecoder struct {
	AbstractDecoder
	reader *bufio.Reader
	State  int
}

func (i IntDiffDecoder) Read() int {
	i.CheckDisposed()
	_, value, _ := lib0.ReadVarInt(i.reader)
	i.State += int(value)
	return int(i.State)
}
