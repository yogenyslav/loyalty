package controller

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_getBackoffSeconds(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name        string
		baseBackoff int
		attempt     int
		maxBackoff  int
		want        int
	}{
		{
			name:        "base backoff 1, attempt 0",
			baseBackoff: 1,
			attempt:     0,
			maxBackoff:  10,
			want:        1,
		},
		{
			name:        "base backoff 2, attempt 1",
			baseBackoff: 2,
			attempt:     1,
			maxBackoff:  10,
			want:        5,
		},
		{
			name:        "base backoff 3, attempt 2",
			baseBackoff: 3,
			attempt:     2,
			maxBackoff:  15,
			want:        15,
		},
		{
			name:        "base backoff exceeds max",
			baseBackoff: 10000,
			attempt:     3,
			maxBackoff:  50,
			want:        50,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			backoff := getBackoffSeconds(tt.baseBackoff, tt.attempt)
			if backoff > tt.maxBackoff {
				backoff = tt.maxBackoff
			}
			assert.InEpsilon(t, tt.want, backoff, 1) // because of jitter
		})
	}
}
