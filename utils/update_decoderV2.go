package utils

import (
	"bufio"
	"io"
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
}

func (u *UpdateDecoderV2) ReadLeftId() *ID {
	// TODO implement me
	panic("implement me")
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

func NewUpdateDecoderV2(reader io.Reader) *UpdateDecoderV2 {
	return &UpdateDecoderV2{}
}
