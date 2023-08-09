package utils

import (
	"bufio"
	"io"

	"YJS-GO/lib0/decoder"
)

var a IDSDecoder = (*DSDecoderV2)(nil)

type DSDecoderV2 struct {
	reader *bufio.Reader
}

func (d *DSDecoderV2) Reader() *bufio.Reader {
	return d.reader
}

func (d *DSDecoderV2) ReadDsLength() uint64 {
	// TODO implement me
	panic("implement me")
}

func (d *DSDecoderV2) ReadDsClock() uint64 {
	// TODO implement me
	panic("implement me")
}

func (d *DSDecoderV2) ResetDsCurVal() {
	// TODO implement me
	panic("implement me")
}

func NewDsDecoderV2(reader io.Reader) *DSDecoderV2 {
	return &DSDecoderV2{reader: bufio.NewReader(reader)}
}

var _ IUpdateDecoder = (*UpdateDecoderV2)(nil)

type UpdateDecoderV2 struct {
	DSDecoderV2
	Keys              []string
	KeyClockDecoder   decoder.IntDiffOptRleDecoder
	ClientDecoder     decoder.UintOptRleDecoder
	LeftClockDecoder  decoder.IntDiffOptRleDecoder
	RightClockDecoder decoder.IntDiffOptRleDecoder
	InfoDecoder       decoder.RleDecoder
	StringDecoder     decoder.StringDecoder
	ParentInfoDecoder decoder.RleDecoder
	TypeRefDecoder    decoder.UintOptRleDecoder
	LengthDecoder     decoder.UintOptRleDecoder
}

func NewUpdateDecoderV2(reader io.Reader) *UpdateDecoderV2 {
	a := bufio.NewReader(reader)
	// Read feature flag - currently unused.
	_, err := a.ReadByte()
	if err != nil {
		panic("get UpdateDecoderV2 panic")
	}
	return &UpdateDecoderV2{
		DSDecoderV2:       DSDecoderV2{reader: a},
		Keys:              []string{},
		KeyClockDecoder:   decoder.IntDiffOptRleDecoder{Reader: decoder.ReadVarUint8ArrayAsStream(a)},
		ClientDecoder:     decoder.UintOptRleDecoder{Reader: decoder.ReadVarUint8ArrayAsStream(a)},
		LeftClockDecoder:  decoder.IntDiffOptRleDecoder{Reader: decoder.ReadVarUint8ArrayAsStream(a)},
		RightClockDecoder: decoder.IntDiffOptRleDecoder{Reader: decoder.ReadVarUint8ArrayAsStream(a)},
		InfoDecoder:       decoder.RleDecoder{Reader: decoder.ReadVarUint8ArrayAsStream(a)},
		StringDecoder:     decoder.StringDecoder{Reader: decoder.ReadVarUint8ArrayAsStream(a)},
		ParentInfoDecoder: decoder.RleDecoder{Reader: decoder.ReadVarUint8ArrayAsStream(a)},
		TypeRefDecoder:    decoder.UintOptRleDecoder{Reader: decoder.ReadVarUint8ArrayAsStream(a)},
		LengthDecoder:     decoder.UintOptRleDecoder{Reader: decoder.ReadVarUint8ArrayAsStream(a)},
	}
}

func (u *UpdateDecoderV2) ReadLeftId() *ID {
	return &ID{u.ClientDecoder.Read(), u.LeftClockDecoder.Read()}
}

func (u *UpdateDecoderV2) ReadRightId() *ID {
	// TODO implement me
	panic("implement me")
}

func (u *UpdateDecoderV2) ReadClient() uint64 {
	// TODO implement me
	panic("implement me")
}

func (u *UpdateDecoderV2) ReadInfo() byte {
	// TODO implement me
	panic("implement me")
}

func (u *UpdateDecoderV2) ReadString() string {
	// TODO implement me
	panic("implement me")
}

func (u *UpdateDecoderV2) ReadParentInfo() bool {
	// TODO implement me
	panic("implement me")
}

func (u *UpdateDecoderV2) ReadTypeRef() uint {
	// TODO implement me
	panic("implement me")
}

func (u *UpdateDecoderV2) ReadLength() uint64 {
	// TODO implement me
	panic("implement me")
}

func (u *UpdateDecoderV2) ReadAny() any {
	// TODO implement me
	panic("implement me")
}

func (u *UpdateDecoderV2) ReadBuffer() []byte {
	// TODO implement me
	panic("implement me")
}

func (u *UpdateDecoderV2) ReadKey() string {
	// TODO implement me
	panic("implement me")
}

func (u *UpdateDecoderV2) ReadJson() any {
	// TODO implement me
	panic("implement me")
}

func (u *UpdateDecoderV2) CheckDisposed() bool {
	return false
}
