package migrations

import "embed"

//go:embed *.sql
var migrations embed.FS

// GetMigrationsFS returns the embedded filesystem containing migration files.
func GetMigrationsFS() embed.FS {
	return migrations
}
