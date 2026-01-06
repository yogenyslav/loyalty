// Package loyalty is the root package for the loyalty application domain-related logic.
package loyalty

import (
	"context"
	"log/slog"

	"github.com/gofiber/fiber/v3"
	"github.com/golang-jwt/jwt/v5"
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
	"github.com/yogenyslav/loyalty/pkg/errs"
)

// JwtProvider defines methods for JWT token management.
type JwtProvider interface {
	CreateAccessToken(userID int64) (string, error)
	ParseAccessToken(accessTokenString string) (*jwt.Token, error)
}

// InitService initializes the loyalty service dependencies.
func InitService(
	ctx context.Context,
	router fiber.Router,
	db database.DB,
	jwtProvider JwtProvider,
	accrual oc.AccrualService,
) {
	mws := user.Middlewares{AuthMiddleware: middleware.WithAuth(jwtProvider)}
	uow := database.NewUnitOfWork(db)

	repos := setupRepos(db)
	controllers := setupControllers(repos, uow, jwtProvider, accrual)
	h := setupHandlers(controllers)

	go startOrdersPolling(ctx, controllers.order)

	userGroup := router.Group("/user")
	user.SetupRoutes(userGroup, mws, h.auth, h.order, h.balance)
}

func startOrdersPolling(ctx context.Context, ctrl *oc.Controller) {
	errCh := make(chan error)
	go ctrl.PollOrders(ctx, errCh)

	for {
		select {
		case <-ctx.Done():
			return
		case err := <-errCh:
			if err != nil {
				slog.Error("order polling error", slog.Any("error", errs.Wrap(err)))
			}
		}
	}
}

type repos struct {
	balance *br.Repo
	auth    *ar.Repo
	order   *or.Repo
}

func setupRepos(db database.DB) *repos {
	return &repos{
		balance: br.New(db),
		auth:    ar.New(db),
		order:   or.New(db),
	}
}

type controllers struct {
	balance *bc.Controller
	auth    *ac.Controller
	order   *oc.Controller
}

func setupControllers(
	r *repos,
	uow database.UnitOfWork,
	jwtProvider JwtProvider,
	accrual oc.AccrualService,
) *controllers {
	return &controllers{
		balance: bc.New(r.balance, uow),
		auth:    ac.New(r.auth, r.balance, jwtProvider, uow),
		order:   oc.New(r.order, r.balance, accrual, uow),
	}
}

type handlers struct {
	balance *bh.Handler
	auth    *ah.Handler
	order   *oh.Handler
}

func setupHandlers(c *controllers) *handlers {
	return &handlers{
		balance: bh.New(c.balance),
		auth:    ah.New(c.auth),
		order:   oh.New(c.order),
	}
}
