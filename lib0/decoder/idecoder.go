package decoder

import (
	"bufio"
	"io"
)

type IDecoder[v any] interface {
	io.Reader
	ReadV() v
}

type BaseDecoder[v any] struct {
	*bufio.Reader
}
