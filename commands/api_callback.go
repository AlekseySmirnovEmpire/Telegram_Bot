package commands

import (
	data2 "Telegram_Bot/data"
	"fmt"
	tgbotapi "github.com/Syfaro/telegram-bot-api"
	"github.com/enescakir/emoji"
	"strings"
)

func callBack(upd *tgbotapi.Update, bot *tgbotapi.BotAPI) (err error) {
	var str string
	data := strings.Split(upd.CallbackQuery.Data, ":")
	if len(data) != 4 {
		return nil
	}

	var showMenu bool
	msg := tgbotapi.NewMessage(upd.CallbackQuery.Message.Chat.ID, "")

	switch data[0] {
	case "Age_Confirm":
		str, showMenu, err = ageAnswerCheck(data)
		break
	case "Survey":
		err = survey(&data, upd, bot)
		return err
	default:
		str = defaultAnswer()
	}

	if str != "" {
		msg.Text = str
		_, _ = bot.Send(msg)
	}
	if showMenu {
		_ = initMainMenu(upd, bot, upd.CallbackQuery.Message.Chat.ID)
	}

	return err
}

func ageAnswerCheck(data []string) (str string, showMenu bool, err error) {
	switch data[1] {
	case "Yes":
		str = fmt.Sprintf("Отлично! Давай приступим!%v", emoji.BeamingFaceWithSmilingEyes)
		err = data2.ChangeAgeConfirm(data[2])
		if err == nil {
			um[data[2]].AgeConfirmed = true
			showMenu = true
		}
		break
	case "No":
		str = fmt.Sprintf("Вам должно быть 18+ для пользования ботом!%v", emoji.FaceWithRollingEyes)
		showMenu = false
		break
	}

	if err != nil {
		return "", false, err
	}
	return str, showMenu, nil
}

func defaultAnswer() string {
	str := fmt.Sprintf(
		"Прости, пока рыбов не продаём, просто показываем, кросивое, да? %v",
		emoji.BeamingFaceWithSmilingEyes)
	return str
}

func removeInlineBlock(chatId int64, msgId int, bot *tgbotapi.BotAPI) {
	editedMsg := tgbotapi.NewEditMessageReplyMarkup(
		chatId,
		msgId,
		tgbotapi.InlineKeyboardMarkup{InlineKeyboard: make([][]tgbotapi.InlineKeyboardButton, 0, 0)})
	_, _ = bot.Send(editedMsg)
}
