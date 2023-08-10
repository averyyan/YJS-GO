package encoder

import (
	"strings"
)

var _ IEncoder[string] = (*StringEncoder)(nil)

type StringEncoder struct {
	Sb            *strings.Builder
	LengthEncoder *UintOptRleEncoder
	Disposed      bool
}

func (s *StringEncoder) Write(v any) {
	s.Sb.WriteString(v.(string))
	a := []rune(v.(string))
	s.LengthEncoder.Write(uint(len(a)))
}

func (s *StringEncoder) Write2(value []rune, offset, count int) {
	s.Sb.WriteString(string(value))
	s.LengthEncoder.Write(uint(count))
}

func (s *StringEncoder) ToArray() []byte {
	return []byte(s.Sb.String())
}

func (s *StringEncoder) GetBuffer() ([]byte, int) {
	panic("{nameof(StringEncoder)} doesn't use temporary byte buffers")
}

func (s *StringEncoder) Dispose(disposing bool) {
	if !s.Disposed {
		if disposing {
			s.Sb.Reset()
			s.LengthEncoder.Dispose()
		}

		s.Sb = nil
		s.LengthEncoder = nil
		s.Disposed = true
	}
}

func (s *StringEncoder) CheckDisposed() {
	if s.Disposed {
		// throw new ObjectDisposedException(GetType().ToString())
	}
}
