package decoder

type IDecoder[v int | uint | uint64 | string | byte] interface {
	Read() v
}

type AbstractDecoder struct {
	LeaveOpen bool
	Disposed  bool
}

func (a AbstractDecoder) CheckDisposed() {
	if a.Disposed {
		return
	}
}

func (a AbstractDecoder) HasContent() bool {
	return false
}
