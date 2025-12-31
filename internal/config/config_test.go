package config

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/yogenyslav/loyalty/internal/accrual"
	"github.com/yogenyslav/loyalty/internal/server"
	"github.com/yogenyslav/loyalty/pkg/database"
	"github.com/yogenyslav/loyalty/pkg/jwt"
)

const testConfigData = `server:
  run_address: "localhost:8080"
database:
  uri: "postgres://pguser:pgpass@localhost:5432/dev?sslmode=disable"
jwt:
  secret: "mysecret"
  expire: 3600
accrual:
  accrual_addr: "http://localhost:8080"
`

func TestNew(t *testing.T) {
	tmpDir := os.TempDir()
	tmpFile := tmpDir + "/config_test.yaml"
	err := os.WriteFile(tmpFile, []byte(testConfigData), 0o600)
	require.NoError(t, err)

	tests := []struct {
		name    string
		path    []string
		want    *Config
		wantErr bool
	}{
		{
			name: "load from yaml file",
			path: []string{tmpFile},
			want: &Config{
				Server: server.Config{
					RunAddress:  "localhost:8080",
					BodyLimitMb: -1,
				},
				DB: database.Config{
					URI: "postgres://pguser:pgpass@localhost:5432/dev?sslmode=disable",
				},
				Jwt: jwt.Config{
					Secret: "mysecret",
					Expire: 3600,
				},
				Accrual: accrual.Config{
					BaseURL:      "http://localhost:8080",
					PollInterval: 10,
					NumWorkers:   5,
				},
			},
			wantErr: false,
		},
		{
			name:    "file does not exist",
			path:    []string{tmpDir + "/asdf.yaml"},
			want:    nil,
			wantErr: true,
		},
		{
			name: "load from env",
			want: &Config{
				Server: server.Config{
					RunAddress:  "localhost:9999",
					BodyLimitMb: -1,
				},
				DB: database.Config{
					URI: "postgres://pguser:pgpass@localhost:5432/dev?sslmode=disable",
				},
				Jwt: jwt.Config{
					Secret: "secret",
				},
				Accrual: accrual.Config{
					BaseURL:      "localhost:8080",
					PollInterval: 10,
					NumWorkers:   5,
				},
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.path == nil {
				t.Setenv("RUN_ADDRESS", "localhost:9999")
				t.Setenv("DATABASE_URI", "postgres://pguser:pgpass@localhost:5432/dev?sslmode=disable")
				t.Setenv("ACCRUAL_SYSTEM_ADDRESS", "localhost:8080")
				t.Setenv("JWT_SECRET", "secret")
			}

			cfg, err := New(tt.path...)
			if tt.wantErr {
				require.Error(t, err)
				assert.Nil(t, cfg)
			} else {
				require.NoError(t, err)
				require.NotNil(t, cfg)
				assert.Equal(t, tt.want, cfg)
			}
		})
	}
}
