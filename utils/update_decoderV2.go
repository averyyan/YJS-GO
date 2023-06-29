package utils

var a IDSDecoder = (*DSDecoderV2)(nil)

type DSDecoderV2 struct {
	IDSDecoder
}

type UpdateDecoderV2 struct {
	DSDecoderV2
	IUpdateDecoder
}
