// Package main is an entry point for the loyalty server application.
package main

import (
	"context"
	"net/http"

	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/pressly/goose/v3"
	"github.com/rs/zerolog/log"
	"github.com/yogenyslav/loyalty/internal/accrual"
	"github.com/yogenyslav/loyalty/internal/config"
	"github.com/yogenyslav/loyalty/internal/loyalty"
	"github.com/yogenyslav/loyalty/internal/server"
	"github.com/yogenyslav/loyalty/migrations"
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
	log.Debug().Interface("config", cfg).Msg("config loaded")

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	pg, err := database.NewPostgres(ctx, cfg.DB.DSN())
	if err != nil {
		return errs.Wrap(err, "connect to db")
	}
	defer pg.Close()

	pgConn, err := pg.SQLDB()
	if err != nil {
		return errs.Wrap(err, "get sql db")
	}
	defer pgConn.Close()

	goose.SetBaseFS(migrations.GetMigrationsFS())
	if err = goose.SetDialect("postgres"); err != nil {
		log.Fatal().Err(err).Msg("set goose dialect")
	}
	err = goose.Up(pgConn, ".")
	if err != nil {
		return errs.Wrap(err, "apply migrations")
	}

	if err = pg.Ping(ctx); err != nil {
		return errs.Wrap(err, "ping db")
	}

	jwtProvider := jwt.New(&cfg.Jwt)
	accrualService := accrual.NewClient(&cfg.Accrual, http.DefaultClient)
	srv := server.New(&cfg.Server)

	apiRouter := srv.Router("/api")
	loyalty.SetupRoutes(ctx, apiRouter, pg, jwtProvider, accrualService)

	if err = srv.Run(); err != nil {
		return errs.Wrap(err, "run server")
	}

	return nil
}
