package middleware

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gofiber/fiber/v3"
	jwtwire "github.com/golang-jwt/jwt/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/yogenyslav/loyalty/internal/server"
)

type mockJwtProvider struct {
	mock.Mock
}

func (m *mockJwtProvider) ParseAccessToken(accessTokenString string) (*jwtwire.Token, error) {
	args := m.Called(accessTokenString)
	return args.Get(0).(*jwtwire.Token), args.Error(1)
}

func testHandler(c fiber.Ctx) error {
	return c.SendStatus(http.StatusOK)
}

func getApp(t *testing.T) *fiber.App {
	t.Helper()
	app := fiber.New(fiber.Config{
		ErrorHandler: server.NewErrorHandler().Handler,
	})
	return app
}

func getReq(t *testing.T) *http.Request {
	t.Helper()
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	return req
}

func TestWithAuth(t *testing.T) {
	t.Run("missing authorization header", func(t *testing.T) {
		app := getApp(t)
		mockProvider := new(mockJwtProvider)
		app.Get("/", WithAuth(mockProvider), testHandler)

		req := getReq(t)
		resp, _ := app.Test(req)
		defer resp.Body.Close()

		assert.Equal(t, http.StatusUnauthorized, resp.StatusCode)
	})

	t.Run("invalid token", func(t *testing.T) {
		app := getApp(t)
		mockProvider := new(mockJwtProvider)
		mockProvider.On("ParseAccessToken", "invalid_token").Return((*jwtwire.Token)(nil), errors.New("invalid token"))
		app.Get("/", WithAuth(mockProvider), testHandler)

		req := getReq(t)
		req.Header.Set("Authorization", "Bearer invalid_token")
		resp, _ := app.Test(req)
		defer resp.Body.Close()

		assert.Equal(t, http.StatusInternalServerError, resp.StatusCode)
	})

	t.Run("valid token", func(t *testing.T) {
		app := getApp(t)
		mockProvider := new(mockJwtProvider)
		token := &jwtwire.Token{
			Claims: jwtwire.MapClaims{"sub": float64(123)},
		}
		mockProvider.On("ParseAccessToken", "valid_token").Return(token, nil)
		app.Get("/", WithAuth(mockProvider), testHandler)

		req := getReq(t)
		req.Header.Set("Authorization", "Bearer valid_token")
		resp, _ := app.Test(req)
		defer resp.Body.Close()

		assert.Equal(t, http.StatusOK, resp.StatusCode)
	})

	t.Run("missing user id in token claims", func(t *testing.T) {
		app := getApp(t)
		mockProvider := new(mockJwtProvider)
		token := &jwtwire.Token{
			Claims: jwtwire.MapClaims{},
		}
		mockProvider.On("ParseAccessToken", "token_without_sub").Return(token, nil)
		app.Get("/", WithAuth(mockProvider), testHandler)

		req := getReq(t)
		req.Header.Set("Authorization", "Bearer token_without_sub")
		resp, _ := app.Test(req)
		defer resp.Body.Close()

		assert.Equal(t, http.StatusUnauthorized, resp.StatusCode)
	})
}
