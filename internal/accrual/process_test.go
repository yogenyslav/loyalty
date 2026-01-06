package accrual

import (
	"bytes"
	"net/http"
	"testing"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

type nopCloser struct {
	*bytes.Buffer
}

func (n nopCloser) Close() error {
	return nil
}

type mockHTTPClient struct {
	mock.Mock
}

func (m *mockHTTPClient) Do(req *http.Request) (*http.Response, error) {
	args := m.Called(req)
	return args.Get(0).(*http.Response), args.Error(1)
}

func TestAccrualClient_ProcessOrder(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		orderNum string
		wantErr  bool
	}{
		{
			name:     "Valid order number",
			orderNum: "2844830162",
			wantErr:  false,
		},
		{
			name:     "Invalid order number",
			orderNum: "12345",
			wantErr:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			mockClient := new(mockHTTPClient)
			if !tt.wantErr {
				mockClient.On("Do", mock.Anything).Return(&http.Response{
					StatusCode: http.StatusOK,
					Body: func() *nopCloser {
						data := []byte(`{"order":"` + tt.orderNum + `","status":"PROCESSED","accrual":100.0}`)
						return &nopCloser{
							Buffer: bytes.NewBuffer(data),
						}
					}(),
				}, nil)
			} else {
				mockClient.On("Do", mock.Anything).Return(&http.Response{
					StatusCode: http.StatusBadRequest,
					Body: func() *nopCloser {
						return &nopCloser{
							Buffer: bytes.NewBuffer([]byte(`{"error":"invalid order number"}`)),
						}
					}(),
				}, nil).Times(3)
			}

			ac := NewService(&Config{
				PollInterval: 1,
				NumWorkers:   5,
			}, mockClient)

			_, err := ac.ProcessOrder(t.Context(), tt.orderNum)
			if tt.wantErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}
		})
	}

	t.Run("HTTP client error", func(t *testing.T) {
		t.Parallel()

		mockClient := new(mockHTTPClient)
		mockClient.On("Do", mock.Anything).Return(&http.Response{}, http.ErrHandlerTimeout).Times(3)

		ac := NewService(&Config{
			PollInterval: 1,
			NumWorkers:   5,
		}, mockClient)

		_, err := ac.ProcessOrder(t.Context(), "2844830162")
		require.Error(t, err)
	})
}
