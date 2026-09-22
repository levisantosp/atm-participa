package db

import (
	"context"
	"database/sql"
	"log"

	"entgo.io/ent/dialect"
	entsql "entgo.io/ent/dialect/sql"
	"github.com/levisantosp/atm-participa/api/ent/generated"
	"github.com/levisantosp/atm-participa/api/utils"

	_ "github.com/jackc/pgx/v5/stdlib"
)

var Client *generated.Client

func Connect() *sql.DB {
	db, err := sql.Open("pgx", utils.Env.DatabaseURL)
	if err != nil {
		log.Fatal(err)
	}

	if err := db.Ping(); err != nil {
		log.Fatal(err)
	}

	Client = generated.NewClient(
		generated.Driver(entsql.OpenDB(dialect.Postgres, db)),
	)

	return db
}

func WithTx[T any](
	ctx context.Context,
	fn func(tx *generated.Tx) (*T, error),
) (*T, error) {
	tx, err := Client.Tx(ctx)
	if err != nil {
		return nil, err
	}

	defer func() {
		_ = tx.Rollback()
	}()

	result, err := fn(tx)
	if err != nil {
		return nil, err
	}

	if err := tx.Commit(); err != nil {
		return nil, err
	}

	return result, nil
}
