package database

import (
	"database/sql"
	"os"

	"github.com/uptrace/bun"
	"github.com/uptrace/bun/dialect/sqlitedialect"
	"github.com/uptrace/bun/driver/sqliteshim"
	"github.com/uptrace/bun/extra/bundebug"
)

func NewSqliteDB() (*bun.DB, error) {
	conn, err := sql.Open(sqliteshim.ShimName, os.Getenv("DATABASE_CONN"))
	if err != nil {
		return nil, err
	}

	db := bun.NewDB(conn, sqlitedialect.New())
	db.AddQueryHook(bundebug.NewQueryHook(
		bundebug.WithVerbose(true),
		bundebug.FromEnv("BUNDEBUG"),
	))

	return db, nil
}
