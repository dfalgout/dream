package migrate

import (
	"database/sql"
	"embed"

	migrate "github.com/rubenv/sql-migrate"
)

//go:embed migrations/*
var dbMigrations embed.FS

type Migration struct {
	migrations migrate.EmbedFileSystemMigrationSource
}

func NewMigration() *Migration {
	migrations := migrate.EmbedFileSystemMigrationSource{
		FileSystem: dbMigrations,
		Root:       "migrations",
	}
	return &Migration{
		migrations,
	}
}

func (m *Migration) Up(conn *sql.DB) (int, error) {
	n, err := migrate.Exec(conn, "sqlite3", m.migrations, migrate.Up)
	if err != nil {
		return 0, err
	}
	return n, nil
}
