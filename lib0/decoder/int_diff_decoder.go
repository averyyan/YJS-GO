package decoder

var _ IDecoder[int] = (*IntDiffDecoder)(nil)

type IntDiffDecoder struct {
}

func (i IntDiffDecoder) Read(p []byte) (n int, err error) {
	// TODO implement me
	panic("implement me")
}

func (i IntDiffDecoder) ReadV() int {
	// TODO implement me
	panic("implement me")
}
