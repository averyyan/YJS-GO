package utils

type YDoc struct {
	GC bool
	GCFilter
}

type ydocOptions struct {
}

func (d YDoc) Clone() *YDoc {
	return &YDoc{}
}

func (d YDoc) EncodeStateVectorV2() []byte {
	var encoder = new(IDSEncoderV2)
}

func WriteStateVector(encoder IDSEncoder) {
	EncodingUtils.WriteStateVector(encoder, Store.GetStateVector())
}
