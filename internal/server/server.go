// Package server implements the loyalty HTTP server.
package server

import (
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/logger"
	recoverer "github.com/gofiber/fiber/v3/middleware/recover"
	"github.com/yogenyslav/loyalty/pkg/errs"
)

// Server represents the loyalty HTTP server.
type Server struct {
	cfg *Config
	app *fiber.App
}

// New creates a new Server instance.
func New(cfg *Config) *Server {
	app := fiber.New(fiber.Config{
		BodyLimit:    cfg.GetBodyLimit(),
		ErrorHandler: NewErrorHandler().Handler,
		AppName:      "LoyaltyProgram API",
	})

	app.Use(logger.New())
	app.Use(recoverer.New())

	return &Server{
		cfg: cfg,
		app: app,
	}
}

// Router returns the server's router with specified prefix.
func (s *Server) Router(prefix string) fiber.Router {
	return s.app.Group(prefix)
}

// Run the HTTP server.
func (s *Server) Run() error {
	errCh := make(chan error, 1)
	stopCh := make(chan os.Signal, 1)
	signal.Notify(stopCh, syscall.SIGINT, syscall.SIGTERM, syscall.SIGABRT)

	go s.listen(errCh)

	select {
	case err := <-errCh:
		return errs.Wrap(err, "server error")
	case <-stopCh:
		slog.Info("shutting down server")
		if err := s.app.Shutdown(); err != nil {
			return errs.Wrap(err, "shutdown server")
		}
	}

	return nil
}

func (s *Server) listen(errCh chan<- error) {
	addr := s.cfg.GetAddr()
	slog.Info("starting server", slog.String("addr", addr))
	if err := s.app.Listen(addr); err != nil {
		errCh <- errs.Wrap(err, "serve http")
	}
}
