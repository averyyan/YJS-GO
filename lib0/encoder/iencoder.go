package encoder

import "io"

type IEncoder[v any] interface {
	Write(v)
	*io.Writer
}
