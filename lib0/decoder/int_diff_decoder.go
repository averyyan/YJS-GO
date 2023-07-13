package decoder

var _ IDecoder[int] = (*IntDiffDecoder)(nil)

// IntDiffDecoder Basic diff encoder using variable length encoding.
// Encodes the values <c>[3, 1100, 1101, 1050, 0]</c> to <c>[3, 1097, 1, -51, -1050]</c>.
// <seealso cref="IntDiffDecoder"/>
type IntDiffDecoder struct {
	State uint64
}

func (i IntDiffDecoder) Read(p []byte) (n int, err error) {
	// TODO implement me
	panic("implement me")
}

func (i IntDiffDecoder) ReadV() int {

	i.reader.WriteVarInt(value - _state)
	_state = value
}
