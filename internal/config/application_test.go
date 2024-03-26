package config

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestApplication_MajorVersion(t *testing.T) {
	tests := []struct {
		name string
		in   string
		exp  int
	}{
		{
			name: "empty",
			in:   "",
			exp:  0,
		},
		{
			name: "no v prefix",
			in:   "1.0.0",
			exp:  1,
		},
		{
			name: "no v prefix greater than 1 digit",
			in:   "23.0.0",
			exp:  23,
		},
		{
			name: "valid",
			in:   "v1.0.0",
			exp:  1,
		},
		{
			name: "valid greater than 1 digit",
			in:   "v157.0.0",
			exp:  157,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			out := Application{version: tt.in}.MajorVersion()
			assert.Equal(t, tt.exp, out)
		})
	}
}
