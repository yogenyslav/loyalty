package accrual

import (
	"bytes"
	"net/http"
	"testing"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestAccrualClient_Poll(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		orders  []string
		wantErr bool
	}{
		{
			name:    "Poll order",
			orders:  []string{"2844830162"},
			wantErr: false,
		},
		{
			name:    "Poll no orders",
			orders:  []string{},
			wantErr: false,
		},
		{
			name:    "Poll with errors",
			orders:  []string{"asdf"},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			mockClient := new(mockHTTPClient)

			if tt.wantErr {
				mockClient.On("Do", mock.Anything).Return(&http.Response{
					StatusCode: http.StatusBadRequest,
					Body: func() *nopCloser {
						return &nopCloser{
							Buffer: bytes.NewBuffer([]byte(`{"error":"invalid order number"}`)),
						}
					}(),
				}, nil).Times(3)
			} else {
				mockClient.On("Do", mock.Anything).Return(&http.Response{
					StatusCode: http.StatusOK,
					Body: func() *nopCloser {
						return &nopCloser{Buffer: bytes.NewBuffer([]byte(`{"order":"123","status":"PROCESSED","accrual":100.0}`))}
					}(),
				}, nil)
			}

			ac := NewService(&Config{
				PollInterval: 1,
				NumWorkers:   2,
			}, mockClient)

			orders := make(chan string, len(tt.orders))
			result := make(chan PollResult, len(tt.orders))

			go ac.Poll(t.Context(), orders, result)

			for _, orderNum := range tt.orders {
				orders <- orderNum
			}
			close(orders)

			var gotErr bool
			for range tt.orders {
				res := <-result
				if res.Err != nil {
					gotErr = true
				}
			}

			require.Equal(t, tt.wantErr, gotErr)
		})
	}
}
