package commands

import (
	"Telegram_Bot/data"
	"fmt"
	tgbotapi "github.com/Syfaro/telegram-bot-api"
	"log"
)

func start(bot *tgbotapi.BotAPI, upd *tgbotapi.Update) error {
	msg := tgbotapi.NewMessage(upd.Message.Chat.ID, upd.Message.Text)

	usr, err := FindUserInArray(int64(upd.Message.From.ID))
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
	if _, err = bot.Send(msg); err != nil {
		return err
	}

	if !CheckAge(int64(upd.Message.From.ID)) {
		msg.Text = fmt.Sprintf(
			"Для продолжения пользования ботом Вы должны быть старше 18 лет.\nВы подтверждаете, что вам больше 18 лет?")
		nkAge := tgbotapi.NewInlineKeyboardMarkup(
			tgbotapi.NewInlineKeyboardRow(
				tgbotapi.NewInlineKeyboardButtonData("Да", fmt.Sprintf("Age_Confirm_Yes:%o", upd.Message.From.ID)),
				tgbotapi.NewInlineKeyboardButtonData("Нет", fmt.Sprintf("Age_Confirm_No:%o", upd.Message.From.ID)),
			))
		msg.ReplyMarkup = nkAge
		if _, err = bot.Send(msg); err != nil {
			return err
		}
	}

	return nil
}
