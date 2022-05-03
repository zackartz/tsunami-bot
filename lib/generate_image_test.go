package lib

import "testing"

func Test_generateBingoImage(t *testing.T) {
	tests := []struct {
		name string
	}{
		{
			name: "hello",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			generateBingoImage()
		})
	}
}
