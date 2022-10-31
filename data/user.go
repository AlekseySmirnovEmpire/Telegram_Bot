package data

import (
	"Telegram_Bot/db"
	"Telegram_Bot/errors"
	errors2 "errors"
	"fmt"
	"github.com/google/uuid"
	"github.com/lib/pq"
	"log"
	"time"
)

type User struct {
	ID             uuid.UUID   `db:"u_id"`
	Key            int64       `db:"user_key"`
	Subscribe      bool        `db:"subscribe"`
	CreatedAt      time.Time   `db:"created_at"`
	SubscribeAt    pq.NullTime `db:"subscribe_at"`
	SubscribeEndAt pq.NullTime `db:"subscribe_end_at"`
}

func InitUser(key int64) (*User, error) {
	u, err := FindUser(key)

	var nse errors.NotSingle
	if err != nil {
		if errors2.As(err, &nse) {
			u, err = createUser(key)
			if err != nil {
				return nil, errors.NotSingle{Val: "users", Err: "cannot create user!"}
			}
		} else {
			return nil, errors.NoConnection{Val: "Postgres", Err: err.Error()}
		}
	}

	return u, nil
}

func FindUser(key int64) (*User, error) {
	query := fmt.Sprintf(`SELECT * FROM users AS u WHERE u.user_key = %o`, key)
	usr, err := db.Select[User](&query)
	if err != nil {
		log.Println(err.Error())
		return nil, err
	}
	if len(usr) != 1 {
		return nil, errors.NotSingle{Val: "users", Err: "There is more then 1 or no data in DB."}
	}

	return usr[0], nil
}

func createUser(key int64) (*User, error) {
	query := fmt.Sprintf(
		`INSERT INTO users (user_key, created_at) VALUES (%o, '%s')`,
		key,
		time.Now().Format("2006-01-02 15:04:05"))
	err := db.InsertOrUpdate[User](&query)
	if err != nil {
		return nil, err
	}

	u, err := FindUser(key)
	if err != nil {
		return nil, err
	}

	return u, nil
}
