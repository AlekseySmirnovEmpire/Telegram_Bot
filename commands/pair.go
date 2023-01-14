package commands

import (
	"fmt"
	tgbotapi "github.com/Syfaro/telegram-bot-api"
	"github.com/enescakir/emoji"
)

func pair(data *[]string,
	upd *tgbotapi.Update,
	bot *tgbotapi.BotAPI) (err error) {

	msg := tgbotapi.NewMessage(upd.CallbackQuery.Message.Chat.ID, "")

	switch (*data)[1] {
	case "Init":
		clearMessagesList(&(*data)[2], upd.CallbackQuery.Message.Chat.ID, bot)
		msg.Text, err = initPair(&msg)
		break
	default:
		msg.Text, err = recognize(data, &msg, bot, upd)
	}

	if err != nil {
		return err
	}

	if msg.Text != "" {
		err = editAndSendMessage(upd.CallbackQuery.Message.Chat.ID, bot, &(*data)[2], &msg)
	}

	return nil
}

func initPair(msg *tgbotapi.MessageConfig) (str string, err error) {
	msg.ReplyMarkup = tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonSwitch(
				"Выбрать",
				"парный опрос")))

	return fmt.Sprintf("Выберите пару из списка своих контактов %s", emoji.SmilingFaceWithSmilingEyes), err
}
