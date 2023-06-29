package decoder

var _ IDecoder[int] = (*IntDiffDecoder)(nil)

type IntDiffDecoder struct {
}
