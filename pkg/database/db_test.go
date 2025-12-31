package database

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestConfig_DSN(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		cfg  *Config
		want string
	}{
		{
			name: "build dsn",
			cfg: &Config{
				Host:     "localhost",
				Port:     "5432",
				User:     "pguser",
				Password: "pgpassword",
				Name:     "dev",
				SSLMode:  "disable",
				Driver:   "postgres",
			},
			want: "postgres://pguser:pgpassword@localhost:5432/dev?sslmode=disable",
		},
		{
			name: "return URI",
			cfg: &Config{
				URI: "postgres://pguser:pgpassword@localhost:5432/dev?sslmode=disable",
			},
			want: "postgres://pguser:pgpassword@localhost:5432/dev?sslmode=disable",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			dsn := tt.cfg.DSN()
			assert.Equal(t, tt.want, dsn)
		})
	}
}
