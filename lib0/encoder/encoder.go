package encoder

import "io"

type binaryEncoder struct {
	io.ByteWriter
}

func (*binaryEncoder) WriteVarUnit(v int) {
	// TODO
}
