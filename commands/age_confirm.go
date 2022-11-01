package commands

import (
	"fmt"
	tgbotapi "github.com/Syfaro/telegram-bot-api"
	"strconv"
)

func ageConfirm(bot *tgbotapi.BotAPI, upd *tgbotapi.Update, strID *string, ans *string) error {
	ID, err := strconv.ParseUint(*strID, 10, 64)
	_, err = FindUserInArray(int64(ID))
	if err != nil {
		return err
	}

	if CheckAge(int64(ID)) {
		return nil
	}

	var answer string

	if *ans == "Age_Confirm_Yes" {
		ageConf[int64(ID)] = true
		answer = fmt.Sprintf("Отлично! Приступим!")
	} else {
		delete(users, int64(ID))
		answer = fmt.Sprintf("Простите, но вам должно быть 18+!")
	}

	_, _ = bot.Send(tgbotapi.NewMessage(upd.CallbackQuery.Message.Chat.ID, answer))
	return nil
}
