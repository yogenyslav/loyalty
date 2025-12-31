package server

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestConfig_GetAddr(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		cfg  *Config
		want string
	}{
		{
			name: "run address",
			cfg: &Config{
				RunAddress: "localhost:9090",
			},
			want: "localhost:9090",
		},
		{
			name: "host and port",
			cfg: &Config{
				Addr: "localhost",
				Port: "1234",
			},
			want: "localhost:1234",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			addr := tt.cfg.GetAddr()
			assert.Equal(t, tt.want, addr)
		})
	}
}

func TestConfig_GetBodyLimit(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		cfg  *Config
		want int
	}{
		{
			name: "custom body limit",
			cfg: &Config{
				BodyLimitMb: 10,
			},
			want: 10 * 1024 * 1024,
		},
		{
			name: "default body limit",
			cfg: &Config{
				BodyLimitMb: -1,
			},
			want: 4 * 1024 * 1024,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			limit := tt.cfg.GetBodyLimit()
			assert.Equal(t, tt.want, limit)
		})
	}
}
