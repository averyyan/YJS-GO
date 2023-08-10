package encoder

type IEncoder[v any] interface {
	Write(v any)
}

type AbstractEncoder struct {
	Disposed bool
}

func (a *AbstractEncoder) FlushV() {

}

func (a *AbstractEncoder) CheckDisposed() {

}

func (a *AbstractEncoder) Dispose(disposing bool) {
	if !a.Disposed {
		if disposing {
			// Stream?.Dispose()
		}

		// Stream = null
		a.Disposed = true
	}
}
