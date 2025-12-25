// Package user provides user API routes, models, and services.
package user

import (
	"github.com/gofiber/fiber/v3"
)

type authHandler interface {
	Login(c fiber.Ctx) error
	Register(c fiber.Ctx) error
}

// SetupRoutes sets up loyalty-related HTTP routes.
func SetupRoutes(router fiber.Router, ah authHandler) {
	router.Post("/register", ah.Register)
	router.Post("/login", ah.Login)
}
