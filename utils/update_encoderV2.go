package utils

type DSEncoderV2 struct {
	IDSEncoder
}

type UpdateEncoderV2S struct {
	DSEncoderV2
}

func (v DSEncoderV2) ToArray() []byte {
	var b []byte
	_, err := v.RestWriter().Read(b)
	if err != nil {
		return nil
	}
	return b
}
