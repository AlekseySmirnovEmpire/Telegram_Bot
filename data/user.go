package data

import (
	"Telegram_Bot/db"
	"Telegram_Bot/errors"
	"fmt"
	"github.com/google/uuid"
	"github.com/lib/pq"
	"log"
	"time"
)

type User struct {
	ID             uuid.UUID   `db:"u_id"`
	Key            string      `db:"user_key"`
	Subscribe      bool        `db:"subscribe"`
	AgeConfirmed   bool        `db:"age_confirmed"`
	CreatedAt      time.Time   `db:"created_at"`
	SubscribeAt    pq.NullTime `db:"subscribe_at"`
	SubscribeEndAt pq.NullTime `db:"subscribe_end_at"`
}

func FindUser(key string) (*User, error) {
	query := fmt.Sprintf(`SELECT * FROM users AS u WHERE u.user_key = '%s'`, key)
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

func CreateUser(key string) (*User, error) {
	query := fmt.Sprintf(
		`INSERT INTO users (user_key, created_at) VALUES ('%s', '%s')`,
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

func InitUsers() (map[string]*User, error) {
	query := fmt.Sprintf(`SELECT * FROM users`)
	ul, err := db.Select[User](&query)
	if err != nil {
		return nil, err
	}
	log.Printf("There are %d users gets from DB!\n", len(ul))

	um := make(map[string]*User, 0)
	for _, u := range ul {
		um[u.Key] = u
	}
	log.Printf("There are %d users inited!\n", len(um))

	return um, nil
}
