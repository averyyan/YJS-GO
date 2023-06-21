package decoder

import (
	"errors"
	"io"

	"YJS-GO/lib0"
)

type binaryDecoder struct {
	io.ByteReader
}

func (d *binaryDecoder) ReadVarUint() (error, uint) {
	num := 0
	len := 0

	for true {
		r, err := d.ReadByte()
		if err != nil {
			return err, 0
		}
		num |= (uint(r) & lib0.Bit7) << len
		len += 7

		if uint(r) < lib0.Bit8 {
			return nil, uint(num)
		}

		if len > 35 {
			return errors.New("integer out of range"), 0
		}
	}
	return nil, 0
}
