package commands

import (
	"fmt"
	tgbotapi "github.com/Syfaro/telegram-bot-api"
	"strconv"
)

func AgeConfirm(msg *tgbotapi.MessageConfig, upd *tgbotapi.Update) string {
	userID := strconv.Itoa(upd.Message.From.ID)
	chatID := strconv.FormatInt(msg.BaseChat.ChatID, 10)

	msg.ReplyMarkup = tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData(
				"Да",
				fmt.Sprintf(
					"Age_Confirm:Yes:%s:%s",
					userID,
					chatID)),
			tgbotapi.NewInlineKeyboardButtonData(
				"Нет",
				fmt.Sprintf(
					"Age_Confirm:No:%s:%s",
					userID,
					chatID)),
		))
	return fmt.Sprintf(
		"Для продолжения пользования ботом Вы должны быть старше 18 лет.\nВы подтверждаете, что вам больше 18 лет?")
}
