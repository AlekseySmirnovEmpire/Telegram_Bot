package data

import (
	"Telegram_Bot/db"
	"fmt"
	"log"
)

type Question struct {
	ID   int    `db:"q_id"`
	Text string `db:"question"`
}

func InitQuestions() ([]*Question, error) {
	log.Println("Start loading questions ....")
	query := fmt.Sprintf(`SELECT * FROM questionnaire AS q ORDER BY q.q_id`)
	ql, err := db.Select[Question](&query)
	if err != nil {
		return nil, err
	}
	log.Printf("Finish loading %d questions!", len(ql))
	return ql, err
}

func InsertAnswer(key, answer *string, questID int) error {
	query := fmt.Sprintf(`INSERT INTO users_to_questions (u_key, q_id, answer) VALUES ('%s', %d, '%s')`,
		*key,
		questID,
		*answer)
	u, err := FindUser(*key)
	if err != nil {
		return err
	}
	u.QuestCount = questID

	log.Printf("Start inserting %s question answer #%d ....", u.Name, questID)
	err = db.InsertOrUpdate[Question](&query)
	if err != nil {
		log.Printf("Failed inserting %s question answer #%d!", u.Name, questID)
		return err
	}

	query = fmt.Sprintf(`UPDATE users SET quest_completed = %d WHERE user_key = '%s'`, questID, *key)
	err = db.InsertOrUpdate[Question](&query)
	if err != nil {
		log.Printf("Failed inserting %s question answer #%d!", u.Name, questID)
		return err
	}
	log.Printf("Update successfully %s question answer #%d!", u.Name, questID)

	return nil
}

func DeleteAnswer(ansID int, userKey *string, full bool) error {
	var query string
	if full {
		query = fmt.Sprintf(`DELETE FROM users_to_questions * WHERE u_key = '%s'`, *userKey)
	} else {
		query = fmt.Sprintf(`DELETE FROM users_to_questions * WHERE u_key = '%s' AND q_id = %d`, *userKey, ansID)
	}

	err := db.Delete[Question](&query)
	if err != nil {
		return err
	}
	return nil
}
