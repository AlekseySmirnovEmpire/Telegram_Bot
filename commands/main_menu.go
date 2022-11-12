package commands

import (
	"Telegram_Bot/myErrors"
	"fmt"
	tgbotapi "github.com/Syfaro/telegram-bot-api"
	"strconv"
	"strings"
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

	var quest tgbotapi.InlineKeyboardButton
	if u.QuestCount > 1 && u.QuestCount < len(ql) {
		quest = tgbotapi.NewInlineKeyboardButtonData(
			"Продолжить анкету",
			fmt.Sprintf(
				"Survey:New:%s:%s",
				userKey,
				chat))
	} else {
		quest = tgbotapi.NewInlineKeyboardButtonData(
			"Пройти анкету",
			fmt.Sprintf(
				"Survey:New:%s:%s",
				userKey,
				chat))
	}

	msg.ReplyMarkup = tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			quest,
			tgbotapi.NewInlineKeyboardButtonData(
				"Пройти парный опрос",
				fmt.Sprintf(
					"Survey:Pair:%s:%s",
					userKey,
					chat)),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData(
				"Отредактировать Анкету",
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

	msg.Text = fmt.Sprintf("%s, добро пожаловать в бота! Выберите действие, что вы хотите сделать:", strings.Split(u.Name, " ")[0])
	_, _ = bot.Send(msg)

	return nil
}
