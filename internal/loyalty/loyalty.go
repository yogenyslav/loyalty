// Package loyalty is the root package for the loyalty application domain-related logic.
package loyalty

import (
	"context"

	"github.com/gofiber/fiber/v3"
	"github.com/golang-jwt/jwt/v5"
	"github.com/rs/zerolog/log"
	"github.com/yogenyslav/loyalty/internal/loyalty/middleware"
	"github.com/yogenyslav/loyalty/internal/loyalty/user"
	ac "github.com/yogenyslav/loyalty/internal/loyalty/user/auth/controller"
	ah "github.com/yogenyslav/loyalty/internal/loyalty/user/auth/handler"
	ar "github.com/yogenyslav/loyalty/internal/loyalty/user/auth/repo"
	bc "github.com/yogenyslav/loyalty/internal/loyalty/user/balance/controller"
	bh "github.com/yogenyslav/loyalty/internal/loyalty/user/balance/handler"
	br "github.com/yogenyslav/loyalty/internal/loyalty/user/balance/repo"
	oc "github.com/yogenyslav/loyalty/internal/loyalty/user/orders/controller"
	oh "github.com/yogenyslav/loyalty/internal/loyalty/user/orders/handler"
	or "github.com/yogenyslav/loyalty/internal/loyalty/user/orders/repo"
	"github.com/yogenyslav/loyalty/pkg/database"
)

// JwtProvider defines methods for JWT token management.
type JwtProvider interface {
	CreateAccessToken(userID int64) (string, error)
	ParseAccessToken(accessTokenString string) (*jwt.Token, error)
}

// SetupRoutes sets up loyalty-related HTTP routes.
func SetupRoutes(
	ctx context.Context,
	router fiber.Router,
	db database.DB,
	jwtProvider JwtProvider,
	accrual oc.AccrualService,
) {
	// user routes
	authMw := middleware.WithAuth(jwtProvider)
	mws := user.Middlewares{
		AuthMiddleware: authMw,
	}

	balanceRepo := br.New(db)
	balanceController := bc.New(balanceRepo)
	balanceHandler := bh.New(balanceController)

	authRepo := ar.New(db, balanceRepo)
	authController := ac.New(authRepo, jwtProvider)
	authHandler := ah.New(authController)

	orderRepo := or.New(db)
	orderController := oc.New(orderRepo, balanceRepo, accrual)
	orderHandler := oh.New(orderController)

	errCh := make(chan error)
	go orderController.PollOrders(ctx, errCh)
	go func() {
		for {
			select {
			case <-ctx.Done():
				return
			case err := <-errCh:
				log.Error().Err(err).Msg("order polling error")
			}
		}
	}()

	userGroup := router.Group("/user")
	user.SetupRoutes(userGroup, mws, authHandler, orderHandler, balanceHandler)
}
