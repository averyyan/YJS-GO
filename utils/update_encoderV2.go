package utils

import (
	"bufio"
	"bytes"

	"YJS-GO/lib0/encoder"
)

var _ IDSEncoder = (*DSEncoderV2)(nil)

type DSEncoderV2 struct {
	writer *bufio.Writer
}

func (v DSEncoderV2) ToArray() []byte {
	var b []byte
	_, err := v.RestWriter().Read(b)
	if err != nil {
		return []byte{}
	}
	return b
}
func (v DSEncoderV2) RestWriter() *bufio.ReadWriter {
	// TODO implement me
	panic("implement me")
}

func (v DSEncoderV2) WriteDsLength(length uint64) {
	// TODO implement me
	panic("implement me")
}

func (v DSEncoderV2) WriteDsClock(clock uint64) {
	// TODO implement me
	panic("implement me")
}

func (v DSEncoderV2) ResetDsCurVal() {
	// TODO implement me
	panic("implement me")
}

type UpdateEncoderV2 struct {
	DSEncoderV2
	keyClock uint
	keyMap   map[string]uint

	keyClockEncoder   encoder.IntDiffOptRleEncoder
	clientEncoder     encoder.UintOptRleEncoder
	leftClockEncoder  encoder.IntDiffOptRleEncoder
	rightClockEncoder encoder.IntDiffOptRleEncoder
	infoEncoder       encoder.RleEncoder
	stringEncoder     encoder.StringEncoder
	parentInfoEncoder encoder.RleEncoder
	typeRefEncoder    encoder.UintOptRleEncoder
	lengthEncoder     encoder.UintOptRleEncoder
}

func NewUpdateEncoderV2() *UpdateEncoderV2 {
	return &UpdateEncoderV2{
		DSEncoderV2:       DSEncoderV2{},
		keyClock:          0,
		keyMap:            nil,
		keyClockEncoder:   encoder.IntDiffOptRleEncoder{},
		clientEncoder:     encoder.UintOptRleEncoder{},
		leftClockEncoder:  encoder.IntDiffOptRleEncoder{},
		rightClockEncoder: encoder.IntDiffOptRleEncoder{},
		infoEncoder:       encoder.RleEncoder{},
		stringEncoder:     encoder.StringEncoder{},
		parentInfoEncoder: encoder.RleEncoder{},
		typeRefEncoder:    encoder.UintOptRleEncoder{},
		lengthEncoder:     encoder.UintOptRleEncoder{},
	}
}

func (u UpdateEncoderV2) Writer() *bufio.Writer {
	// TODO implement me
	panic("implement me")
}

func (u UpdateEncoderV2) WriteLeftId(id *ID) {
	// TODO implement me
	panic("implement me")
}

func (u UpdateEncoderV2) WriteRightId(id *ID) {
	// TODO implement me
	panic("implement me")
}

func (u UpdateEncoderV2) WriteClient(client uint64) {
	// TODO implement me
	panic("implement me")
}

func (u UpdateEncoderV2) WriteInfo(info byte) {
	// TODO implement me
	panic("implement me")
}

func (u UpdateEncoderV2) WriteString(s string) {
	// TODO implement me
	panic("implement me")
}

func (u UpdateEncoderV2) WriteParentInfo(isYKey bool) {
	// TODO implement me
	panic("implement me")
}

func (u UpdateEncoderV2) WriteTypeRef(info uint) {
	// TODO implement me
	panic("implement me")
}

func (u UpdateEncoderV2) WriteLength(len int) {
	// TODO implement me
	panic("implement me")
}

func (u UpdateEncoderV2) WriteAny(object any) {
	// TODO implement me
	panic("implement me")
}

func (u UpdateEncoderV2) WriteBuffer(buf []byte) {
	// TODO implement me
	panic("implement me")
}

func (u UpdateEncoderV2) WriteKey(key string) {
	// TODO implement me
	panic("implement me")
}

func (u UpdateEncoderV2) WriteJson(T any) {
	// TODO implement me
	panic("implement me")
}

func NewDsEncoderV2() *DSEncoderV2 {
	return &DSEncoderV2{writer: bufio.NewWriter(bytes.NewBuffer(oneMBytes))}
}
