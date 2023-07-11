package encoder

var _ IEncoder = (*IncUintOptRleEncoder)(nil)

type IncUintOptRleEncoder struct {
}

func (i *IncUintOptRleEncoder[v any]) Write(v v) {
	// TODO implement me
	panic("implement me")
}



