package commands

import (
	"Telegram_Bot/data"
	"fmt"
	tgbotapi "github.com/Syfaro/telegram-bot-api"
	"github.com/enescakir/emoji"
	"strconv"
)

func start(upd *tgbotapi.Update) (string, error) {
	if IsUserAuth(strconv.Itoa(upd.Message.From.ID)) {
		return fmt.Sprintf("Вы уже стартовали %v", emoji.SlightlySmilingFace), nil
	}

	userName := fmt.Sprintf("%s %s", upd.Message.From.FirstName, upd.Message.From.LastName)
	userKey := strconv.Itoa(upd.Message.From.ID)

	usr, err := data.CreateUser(&userKey, &userName, &upd.Message.From.UserName)
	if err != nil {
		return "", err
	}

	addUser(usr)

	return fmt.Sprintf("Здравствуйте, %s! Добро пожаловать!", upd.Message.From.FirstName), nil
}

func getRandomShit(msg *tgbotapi.Message) string {
	var str string

	if msg.Sticker != nil {
		str = fmt.Sprintf(
			"Я люблю стикеры %v, но продолжить работу смогу только после ответа в опроснике %v",
			emoji.SlightlySmilingFace,
			emoji.BackhandIndexPointingDown.Tone(emoji.Light))
	} else if msg.Photo != nil {
		str = fmt.Sprintf(
			"Фотка огонь %v, но для продолжения вам нужно ответить по кнопке %v",
			emoji.SlightlySmilingFace,
			emoji.BackhandIndexPointingDown.Tone(emoji.Light))
	} else if msg.PinnedMessage != nil || msg.ReplyToMessage != nil {
		str = fmt.Sprintf(
			"Уверен, там что-то интересное %v, но продолжить работу смогу только после ответа в опроснике %v",
			emoji.SlightlySmilingFace,
			emoji.BackhandIndexPointingDown.Tone(emoji.Light))
	} else {
		str = fmt.Sprintf(
			"Для того, чтобы воспользоваться ботом, пожалуйста, ответьте в опроснике %v",
			emoji.SlightlySmilingFace)
	}

	return str
}
