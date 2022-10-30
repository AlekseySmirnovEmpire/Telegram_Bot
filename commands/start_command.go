package commands

import tgbotapi "github.com/Syfaro/telegram-bot-api"

var nkStart = tgbotapi.NewReplyKeyboard(
	tgbotapi.NewKeyboardButtonRow(
		tgbotapi.NewKeyboardButton("Пройти тест!"),
		tgbotapi.NewKeyboardButton("Карта, карта, карта..."),
		tgbotapi.NewKeyboardButton("Здравствуй, карта!"),
	),
)

func start(bot *tgbotapi.BotAPI, upd tgbotapi.Update) error {
	msg := tgbotapi.NewMessage(upd.Message.Chat.ID, upd.Message.Text)

	msg.ReplyMarkup = nkStart
	if _, err := bot.Send(msg); err != nil {
		return err
	}

	return nil
}
