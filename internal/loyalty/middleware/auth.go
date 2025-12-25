// Package middleware provides middleware functions for loyalty HTTP requests.
package middleware

import (
	"github.com/gofiber/fiber/v3"
	jwtwire "github.com/golang-jwt/jwt/v5"
	"github.com/yogenyslav/loyalty/pkg/errs"
	"github.com/yogenyslav/loyalty/pkg/jwt"
)

type jwtProvider interface {
	ParseAccessToken(accessTokenString string) (*jwtwire.Token, error)
}

// WithAuth is a middleware to check access tokens for HTTP requests.
func WithAuth(jwtProvider jwtProvider) fiber.Handler {
	return func(c fiber.Ctx) error {
		authorizationHeader := c.Get("Authorization")
		if authorizationHeader == "" || len(authorizationHeader) < len(jwt.TypeBearerToken)+1 {
			return errs.Wrap(errs.ErrMisingAccessToken, "authorization header")
		}

		accessTokenString := authorizationHeader[len(jwt.TypeBearerToken)+1:]
		token, err := jwtProvider.ParseAccessToken(accessTokenString)
		if err != nil {
			return errs.Wrap(err, "parse access token")
		}

		userID, ok := token.Claims.(jwtwire.MapClaims)["sub"].(int64)
		if !ok {
			return errs.Wrap(errs.ErrInvalidToken, "get user id from token claims")
		}

		c.Locals("userID", userID)

		return c.Next()
	}
}
