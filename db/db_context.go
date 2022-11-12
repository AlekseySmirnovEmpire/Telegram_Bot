package db

import (
	"Telegram_Bot/config"
	"Telegram_Bot/myErrors"
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
	log.Println("Closing DB ....")
	_ = db.Close()
	log.Println("DB Closed!")
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
		return myErrors.NoConnection{Val: "Postgres", Key: err.Error(), Err: err}
	}

	if err = db.Ping(); err != nil {
		return myErrors.NoConnection{Val: "Postgres PING", Key: err.Error(), Err: err}
	}

	log.Println("Connecting DB SUCCESS!")
	return nil
}

func Select[T comparable](query *string) ([]*T, error) {
	obj := make([]*T, 0)
	rows, err := db.Queryx(*query)
	if err != nil {
		log.Println(err.Error())
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var o T
		err = rows.StructScan(&o)
		if err != nil {
			log.Println(err.Error())
			continue
		}
		obj = append(obj, &o)
	}

	return obj, nil
}

func InsertOrUpdate[T comparable](query *string) error {
	_, err := db.Exec(*query)
	if err == nil {
		return err
	}
	return nil
}
