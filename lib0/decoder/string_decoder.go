package decoder

import "bufio"

var _ IDecoder[string] = (*StringDecoder)(nil)

type StringDecoder struct {
	AbstractDecoder
	Reader *bufio.Reader

	LengthDecoder UintOptRleDecoder
	Value         string
	Pos           int
	Disposed      bool
}

func (s StringDecoder) Read() string {
	s.CheckDisposed()

	var length = s.LengthDecoder.Read()
	if length == 0 {
		return ""
	}

	var result = s.Value[s.Pos : s.Pos+int(length)]
	s.Pos += int(length)

	// No need to keep the string in memory anymore.
	// This also covers the case when nothing but empty strings are left.
	if s.Pos >= len(s.Value) {
		s.Value = ""
	}

	return result
}
