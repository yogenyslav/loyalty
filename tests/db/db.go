package db_test

import (
	"bytes"
	"os"
	"testing"
	"time"

	"github.com/ilyakaznacheev/cleanenv"
	"github.com/joho/godotenv"
	"github.com/pressly/goose/v3"
	"github.com/rs/zerolog/log"
	"github.com/stretchr/testify/require"
	"github.com/yogenyslav/loyalty/migrations"
	"github.com/yogenyslav/loyalty/pkg/database"

	_ "embed"

	_ "github.com/jackc/pgx/v5/stdlib"
)

//go:embed test.env
var testEnv []byte

var (
	testDB  database.DB
	testCfg database.Config
)

func init() {
	var err error

	envVars, err := godotenv.Parse(bytes.NewReader(testEnv))
	if err != nil {
		log.Fatal().Err(err).Msg("parse test.env file")
	}
	for k, v := range envVars {
		if err = os.Setenv(k, v); err != nil {
			log.Fatal().Err(err).Msgf("set env var %s", k)
		}
	}

	goose.SetBaseFS(migrations.GetMigrationsFS())
	if err = goose.SetDialect("postgres"); err != nil {
		log.Fatal().Err(err).Msg("set goose dialect")
	}

	if err = cleanenv.ReadEnv(&testCfg); err != nil {
		log.Fatal().Err(err).Msg("read test db config")
	}
}

// SetupTestDB initializes the test database and applies migrations.
func SetupTestDB(t *testing.T) database.DB {
	t.Helper()

	var err error

	testDB, err = database.NewPostgres(t.Context(), testCfg.DSN())
	require.NoError(t, err, "create test db connection")

	RunMigrations(t)

	return testDB
}

// RunMigrations applies all migrations to the test database.
func RunMigrations(t *testing.T) {
	t.Helper()

	sqlDB, err := testDB.SQLDB()
	require.NoError(t, err, "get sql db connection")
	defer sqlDB.Close()

	err = goose.Up(sqlDB, ".")
	require.NoError(t, err, "apply migrations to test db")

	time.Sleep(time.Millisecond * 100)
}

// DropMigrations rolls back all applied migrations in the test database.
func DropMigrations(t *testing.T) {
	t.Helper()

	sqlDB, err := testDB.SQLDB()
	require.NoError(t, err, "get sql db connection")
	defer sqlDB.Close()

	err = goose.DownTo(sqlDB, ".", 0)
	require.NoError(t, err, "rollback migrations in test db")

	time.Sleep(time.Millisecond * 100)
}

// ClearTables deletes all data from the specified tables in the test database.
func ClearTables(t *testing.T, tables ...string) {
	t.Helper()

	ctx := t.Context()
	for _, table := range tables {
		_, err := testDB.Exec(ctx, "delete from "+table)
		require.NoError(t, err, "clear table "+table)
	}
}
