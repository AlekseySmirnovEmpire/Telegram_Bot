package commands

import (
	"Telegram_Bot/data"
	"fmt"
	tgbotapi "github.com/Syfaro/telegram-bot-api"
	"log"
)

var (
	nkStart = tgbotapi.NewReplyKeyboard(
		tgbotapi.NewKeyboardButtonRow(
			tgbotapi.NewKeyboardButton("Пройти тест!"),
			tgbotapi.NewKeyboardButton("Карта, карта, карта..."),
			tgbotapi.NewKeyboardButton("Здравствуй, карта!"),
		),
	)
)

func start(bot *tgbotapi.BotAPI, upd *tgbotapi.Update) error {
	msg := tgbotapi.NewMessage(upd.Message.Chat.ID, upd.Message.Text)

	usr, err := data.FindUser(int64(upd.Message.From.ID))
	if err != nil {
		usr, err = data.InitUser(int64(upd.Message.From.ID))
		if err != nil {
			msg.Text = fmt.Sprintf(
				"Здравствуйте, %s!\nПриносим свои извенения, бот недоступен по техническим причинам!",
				upd.Message.From.FirstName)
			log.Println(err.Error())
			_, _ = bot.Send(msg)
			return nil
		}

		AddUser(usr, usr.Key)
	}

	msg.Text = fmt.Sprintf("Здравствуйте, %s! Добро пожаловать!", upd.Message.From.FirstName)
	msg.ReplyMarkup = nkStart
	if _, err = bot.Send(msg); err != nil {
		return err
	}

	return nil
}
