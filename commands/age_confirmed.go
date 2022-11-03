package commands

import (
	"fmt"
	tgbotapi "github.com/Syfaro/telegram-bot-api"
)

func AgeConfirm(msg *tgbotapi.MessageConfig, upd *tgbotapi.Update) string {
	msg.ReplyMarkup = tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("Да", fmt.Sprintf("Age_Confirm_Yes:%o", upd.Message.From.ID)),
			tgbotapi.NewInlineKeyboardButtonData("Нет", fmt.Sprintf("Age_Confirm_No:%o", upd.Message.From.ID)),
		))
	return fmt.Sprintf(
		"Для продолжения пользования ботом Вы должны быть старше 18 лет.\nВы подтверждаете, что вам больше 18 лет?")
}
