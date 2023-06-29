package utils

import "testing"

func TestGenerateRsaKey(t *testing.T) {
	type args struct {
		keySize int
	}
	tests := []struct {
		name string
		args args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			GenerateRsaKey(tt.args.keySize)
		})
	}
}
