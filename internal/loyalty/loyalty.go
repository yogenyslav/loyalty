// Package loyalty is the root package for the loyalty application domain-related logic.
package loyalty

import (
	"github.com/gofiber/fiber/v3"
	"github.com/golang-jwt/jwt/v5"
	"github.com/yogenyslav/loyalty/internal/loyalty/user"
	ac "github.com/yogenyslav/loyalty/internal/loyalty/user/auth/controller"
	ah "github.com/yogenyslav/loyalty/internal/loyalty/user/auth/handler"
	ar "github.com/yogenyslav/loyalty/internal/loyalty/user/auth/repo"
	"github.com/yogenyslav/loyalty/pkg/database"
)

// JwtProvider defines methods for JWT token management.
type JwtProvider interface {
	CreateAccessToken(userID int64) (string, error)
	ParseAccessToken(accessTokenString string) (*jwt.Token, error)
}

// SetupRoutes sets up loyalty-related HTTP routes.
func SetupRoutes(router fiber.Router, db database.DB, jwtProvider JwtProvider) {
	authRepo := ar.New(db)
	authController := ac.New(authRepo, jwtProvider)
	authHandler := ah.New(authController)

	userGroup := router.Group("/user")
	user.SetupRoutes(userGroup, authHandler)
}
