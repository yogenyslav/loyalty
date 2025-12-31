package luhn

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_checksum(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name   string
		number string
		want   int
	}{
		{
			name:   "checksum is calculated",
			number: "1",
			want:   1,
		},
		{
			name:   "checksum is calculated",
			number: "18",
			want:   0,
		},
		{
			name:   "checksum is calculated",
			number: "123450",
			want:   5,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			next := checksum(tt.number)
			assert.Equal(t, tt.want, next)
		})
	}
}

func TestValidate(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name   string
		number string
		want   bool
	}{
		{
			name:   "valid number",
			number: "4532015112830366",
			want:   true,
		},
		{
			name:   "invalid number",
			number: "1234567890123456",
			want:   false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			isValid := Validate(tt.number)
			assert.Equal(t, tt.want, isValid)
		})
	}
}
