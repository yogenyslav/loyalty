package retry_test

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/yogenyslav/ya-metrics/pkg/retry"
)

func TestWithLinearBackoffRetry(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		cfg     *retry.Config
		fn      func(context.Context) error
		wantErr bool
	}{
		{
			name: "Success on first try",
			cfg: &retry.Config{
				MaxRetries:         3,
				LinearBackoffMilli: 100,
			},
			fn: func(context.Context) error {
				return nil
			},
			wantErr: false,
		},
		{
			name: "Success on third try",
			cfg: &retry.Config{
				MaxRetries:         3,
				LinearBackoffMilli: 100,
			},
			fn: func() func(context.Context) error {
				cnt := 0
				return func(context.Context) error {
					if cnt < 2 {
						cnt++
						return errors.New("err")
					}
					return nil
				}
			}(),
			wantErr: false,
		},
		{
			name: "Fail after max retries",
			cfg: &retry.Config{
				MaxRetries:         2,
				LinearBackoffMilli: 100,
			},
			fn: func(context.Context) error {
				return errors.New("err")
			},
			wantErr: true,
		},
		{
			name: "Nil config",
			cfg:  nil,
			fn: func(context.Context) error {
				return nil
			},
			wantErr: false,
		},
		{
			name: "Nil config with error",
			cfg:  nil,
			fn: func(context.Context) error {
				return errors.New("err")
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			err := retry.WithLinearBackoffRetry(context.Background(), tt.cfg, tt.fn)
			if tt.wantErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}
		})
	}
}
