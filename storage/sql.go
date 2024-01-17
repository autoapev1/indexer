package store

import (
	"database/sql"
	"fmt"

	"github.com/autoapev1/indexer/config"
	"github.com/uptrace/bun"
	"github.com/uptrace/bun/dialect/pgdialect"
	"github.com/uptrace/bun/driver/pgdriver"
	"github.com/uptrace/bun/extra/bundebug"
)

type PostgresDB struct {
	DB *bun.DB
}

func NewPostgresDB(conf config.Postgres) *PostgresDB {
	PostgresDB := &PostgresDB{}
	uri := fmt.Sprintf(
		"user=%s password=%s dbname=%s host=%s port=%s sslmode=%s",
		conf.User,
		conf.Password,
		conf.Name,
		conf.Host,
		conf.Port,
		conf.SSLMode,
	)

	pgdb := sql.OpenDB(pgdriver.NewConnector(pgdriver.WithDSN(uri)))

	db := bun.NewDB(pgdb, pgdialect.New())
	PostgresDB.DB = db

	db.AddQueryHook(bundebug.NewQueryHook(bundebug.WithVerbose(true)))

	return PostgresDB
}
