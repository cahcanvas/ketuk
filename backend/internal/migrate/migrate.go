package migrate

import (
	"database/sql"
	"embed"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/stdlib"
	"github.com/pressly/goose/v3"
)

//go:embed sql/*.sql
var files embed.FS

func Up(databaseURL string) error {
	goose.SetBaseFS(files)
	if err := goose.SetDialect("postgres"); err != nil {
		return err
	}

	// Same pooler-compatibility fix as internal/db: disable pgx's
	// server-side prepared statement cache, which collides with poolers
	// (e.g. Supabase's Supavisor in transaction mode) that route each
	// statement to a different backend connection.
	connCfg, err := pgx.ParseConfig(databaseURL)
	if err != nil {
		return fmt.Errorf("parse database url for migrate: %w", err)
	}
	connCfg.DefaultQueryExecMode = pgx.QueryExecModeSimpleProtocol
	dsn := stdlib.RegisterConnConfig(connCfg)
	defer stdlib.UnregisterConnConfig(dsn)

	db, err := sql.Open("pgx", dsn)
	if err != nil {
		return fmt.Errorf("open db for migrate: %w", err)
	}
	defer db.Close()
	return goose.Up(db, "sql")
}
