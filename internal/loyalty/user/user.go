// Package user provides user API routes, models, and services.
package user

import (
	"github.com/gofiber/fiber/v3"
)

type authHandler interface {
	Login(c fiber.Ctx) error
	Register(c fiber.Ctx) error
}

type orderHandler interface {
	ProcessOrder(c fiber.Ctx) error
	ListOrders(c fiber.Ctx) error
}

type balanceHandler interface {
	GetBalance(c fiber.Ctx) error
	Withdraw(c fiber.Ctx) error
	ListWithdrawals(c fiber.Ctx) error
}

// Middlewares stores user-related middlewares.
type Middlewares struct {
	AuthMiddleware fiber.Handler
}

// SetupRoutes sets up loyalty-related HTTP routes.
func SetupRoutes(router fiber.Router, mws Middlewares, ah authHandler, oh orderHandler, bh balanceHandler) {
	// auth routes
	router.Post("/register", ah.Register)
	router.Post("/login", ah.Login)

	// order routes
	router.Post("/orders", mws.AuthMiddleware, oh.ProcessOrder)
	router.Get("/orders", mws.AuthMiddleware, oh.ListOrders)

	// balance routes
	router.Get("/balance", mws.AuthMiddleware, bh.GetBalance)
	router.Post("/balance/withdraw", mws.AuthMiddleware, bh.Withdraw)
	router.Get("/withdrawals", mws.AuthMiddleware, bh.ListWithdrawals)
}
