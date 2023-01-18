package data

import (
	"Telegram_Bot/db"
	"fmt"
	"time"

	"github.com/enescakir/emoji"
	"github.com/google/uuid"
)

type UserQuestionnary struct {
	ID            uuid.UUID `db:"q_id"`
	fileName      string    `db:"file_name"`
	fileExtension string    `db:"file_ext"`
	fileData      []byte    `db:"file_data"`
	CreatedAt     time.Time `db:"created_at"`
}

func ReturnUserQuestinaryId(userId *string, isSingl bool) (str string, err error) {
	query := fmt.Sprintf(`SELECT UQ.q_id FROM user_questionnaire AS UQ 
	JOIN user_to_anket AS UTA ON UTA.q_id = UQ.q_id
	JOIN users AS U ON U.u_id = UTA.u_id
	WHERE UTA.u_id = '%s' AND UTA.is_single = %t`, *userId, isSingl)

	id, err := db.Select[uuid.UUID](&query)
	if err != nil {
		return
	}
	if len(id) == 0 {
		str = fmt.Sprintf("У вас пока что ещё нет личной анкеты %s", emoji.WinkingFace)
		return
	}

	str = fmt.Sprintf("ID вашей анкеты: `%s`\nВ целях безопасности не передавайте его посторонним людям %s", id[0].String(), emoji.WinkingFace)

	return
}

func CreateAnket() (str string, err error) {
	return "Yes", nil
}
