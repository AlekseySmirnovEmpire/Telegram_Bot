package commands

import (
	"Telegram_Bot/myErrors"
	"fmt"
	tgbotapi "github.com/Syfaro/telegram-bot-api"
	"strconv"
)

func initMainMenu(upd *tgbotapi.Update, bot *tgbotapi.BotAPI, chatID int64) (err error) {
	var userKey string
	if upd.Message == nil {
		userKey = strconv.Itoa(upd.CallbackQuery.From.ID)
	} else {
		userKey = strconv.Itoa(upd.Message.From.ID)
	}

	u, err := FindUser(userKey)
	if err != nil {
		return err
	}
	if !u.AgeConfirmed {
		return myErrors.NotConfirmed{Val: "возраст", Key: "меню"}
	}

	msg := tgbotapi.NewMessage(chatID, "")
	chat := strconv.FormatInt(chatID, 10)

	msg.ReplyMarkup = tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData(
				"Пройти опрос",
				fmt.Sprintf(
					"Survey:New:%s:%s",
					userKey,
					chat)),
			tgbotapi.NewInlineKeyboardButtonData(
				"Пройти парный опрос",
				fmt.Sprintf(
					"Survey:Pair:%s:%s",
					userKey,
					chat)),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData(
				"Пройти опрос заново",
				fmt.Sprintf(
					"Survey:Old:%s:%s",
					userKey,
					chat)),
			tgbotapi.NewInlineKeyboardButtonData(
				"Найти пару",
				fmt.Sprintf(
					"Find:Pair:%s:%s",
					userKey,
					chat)),
		))

	msg.Text = fmt.Sprintf("%s, добро пожаловать в бота! Выберите действие, что вы хотите сделать:", u.Name)
	_, _ = bot.Send(msg)

	return nil
}
