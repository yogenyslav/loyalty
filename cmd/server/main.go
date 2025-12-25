// Package main is an entry point for the loyalty server application.
package main

import (
	"context"

	"github.com/rs/zerolog/log"
	"github.com/yogenyslav/loyalty/internal/config"
	"github.com/yogenyslav/loyalty/internal/loyalty"
	"github.com/yogenyslav/loyalty/internal/server"
	"github.com/yogenyslav/loyalty/pkg/database"
	"github.com/yogenyslav/loyalty/pkg/errs"
	"github.com/yogenyslav/loyalty/pkg/jwt"
)

func main() {
	if err := run(); err != nil {
		log.Fatal().Err(err).Msg("fatal error")
	}
}

func run() error {
	cfg, err := config.New()
	if err != nil {
		return errs.Wrap(err, "load config")
	}

	ctx := context.Background()
	pg, err := database.NewPostgres(ctx, cfg.DB.DSN)
	if err != nil {
		return errs.Wrap(err, "connect to db")
	}
	defer pg.Close()

	jwtProvider := jwt.New(cfg.Jwt)
	srv := server.New(cfg.Server)

	apiRouter := srv.Router("/api")
	loyalty.SetupRoutes(apiRouter, pg, jwtProvider)

	if err := srv.Run(); err != nil {
		return errs.Wrap(err, "run server")
	}

	return nil
}
