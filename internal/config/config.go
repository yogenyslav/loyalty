// Package config provides configuration management for the loyalty application.
package config

import (
	"flag"
	"os"

	"github.com/ilyakaznacheev/cleanenv"
	"github.com/yogenyslav/loyalty/internal/server"
	"github.com/yogenyslav/loyalty/pkg/database"
	"github.com/yogenyslav/loyalty/pkg/errs"
	"github.com/yogenyslav/loyalty/pkg/jwt"
)

// Config holds the entire application configuration settings.
type Config struct {
	Server *server.Config   `yaml:"server"`
	DB     *database.Config `yaml:"database"`
	Jwt    *jwt.Config      `yaml:"jwt"`
}

// New creates a new Config.
func New(path ...string) (*Config, error) {
	flags := flag.NewFlagSet("loyalty", flag.ContinueOnError)
	serverAddr := flags.String("a", "", "адрес и порт запуска сервиса")
	dbURI := flags.String("d", "", "адрес подключения к базе данных")
	accrualAddr := flags.String("r", "", "адрес системы расчёта начислений")

	if err := flags.Parse(os.Args[1:]); err != nil {
		return nil, errs.Wrap(err, "parse flags")
	}

	cfg := Config{
		Server: &server.Config{
			RunAddress:  *serverAddr,
			AccrualAddr: *accrualAddr,
		},
		DB: &database.Config{
			DSN: *dbURI,
		},
	}

	if len(path) > 0 {
		if err := cleanenv.ReadConfig(path[0], &cfg); err != nil {
			return nil, errs.Wrap(err, "read yaml config")
		}
	} else {
		if err := cleanenv.ReadEnv(&cfg); err != nil {
			return nil, errs.Wrap(err, "read env config")
		}
	}

	return &cfg, nil
}
