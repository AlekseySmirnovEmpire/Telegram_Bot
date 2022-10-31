package db

import (
	"Telegram_Bot/config"
	"Telegram_Bot/errors"
	"github.com/jmoiron/sqlx"
	_ "github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"log"
)

const dbName = "postgres"

var (
	db *sqlx.DB
)

// CloseDB close DB.
func CloseDB() {
	_ = db.Close()
}

// InitDB init db connection.
func InitDB() error {
	if db != nil {
		return nil
	}

	log.Println("Connecting DB ....")

	var err error
	db, err = sqlx.Open(dbName, config.ConStr)
	if err != nil {
		return errors.NoConnection{Val: "Postgres", Err: err.Error()}
	}

	if err = db.Ping(); err != nil {
		return errors.NoConnection{Val: "Postgres PING", Err: err.Error()}
	}

	log.Println("Connecting DB SUCCESS!")
	return nil
}
