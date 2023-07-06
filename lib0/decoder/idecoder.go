package decoder

import "io"

type IDecoder[v any] interface {
	io.Reader
	ReadV() v
}
