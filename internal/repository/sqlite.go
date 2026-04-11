package repository

import (
	"fmt"

	"github.com/jmoiron/sqlx"
	_ "github.com/glebarez/go-sqlite"
)

const(
	cardsTable = "cards"
	keysTable = "keys"
	terminalsTable = "terminals"
	transactionsTable = "transactions"
	usersTable = "users"
)

type Config struct {
	DBPath string
}

func NewSQLiteDB(cfg Config) (*sqlx.DB, error) {
	db, err := sqlx.Open("sqlite", cfg.DBPath)
	if err != nil {
		return nil, err
	}

	if err = db.Ping(); err != nil {
		return nil, err
	}

	fmt.Println("SQLite", cfg.DBPath, "connected successfully.")
	return db, nil
}
