package content

import "testing"

func TestAny_Copy(t *testing.T) {
	var a = struct {
		b any
	}{
		[]any{},
	}

	a.b = append(a.b.([]any), struct{}{})
}
