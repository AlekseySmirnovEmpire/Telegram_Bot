package commands

import (
	data2 "Telegram_Bot/data"
	"fmt"
	tgbotapi "github.com/Syfaro/telegram-bot-api"
	"github.com/enescakir/emoji"
	"strconv"
	"strings"
)

func callBack(upd *tgbotapi.Update, bot *tgbotapi.BotAPI) (err error) {
	var str string
	data := strings.Split(upd.CallbackQuery.Data, ":")
	if len(data) != 4 {
		return nil
	}

	switch data[0] {
	case "Age_Confirm":
		str, err = ageAnswerCheck(bot, data)
		break
	default:
		str = defaultAnswer(bot, data)
	}

	if str != "" {
		_, _ = bot.Send(tgbotapi.NewMessage(upd.CallbackQuery.Message.Chat.ID, str))
	}

	return err
}

func ageAnswerCheck(bot *tgbotapi.BotAPI, data []string) (str string, err error) {
	ch := make(chan struct{})
	defer close(ch)
	go func() {
		chatId, _ := strconv.ParseInt(data[3], 10, 64)
		msgId := ageConfMesID[data[2]]
		removeInlineBlock(chatId, msgId, bot)
		ch <- struct{}{}
	}()

	switch data[1] {
	case "Yes":
		str = fmt.Sprintf("Отлично! Давай приступим!%v", emoji.BeamingFaceWithSmilingEyes)
		err = data2.ChangeAgeConfirm(data[2])
		if err == nil {
			um[data[2]].AgeConfirmed = true
			if _, ok := ageConfMesID[data[2]]; ok {
				delete(ageConfMesID, data[2])
			}
		}
		break
	case "No":
		str = fmt.Sprintf("Вам должно быть 18+ для пользования ботом!%v", emoji.FaceWithRollingEyes)
		if _, ok := ageConfMesID[data[2]]; ok {
			delete(ageConfMesID, data[2])
		}
		break
	}

	<-ch
	if err != nil {
		return "", err
	}
	return str, nil
}

func defaultAnswer(bot *tgbotapi.BotAPI, data []string) string {
	ch := make(chan struct{})
	defer close(ch)
	go func() {
		chatId, _ := strconv.ParseInt(data[3], 10, 64)
		msgId, _ := strconv.Atoi(data[4])
		removeInlineBlock(chatId, msgId, bot)
		ch <- struct{}{}
	}()

	str := fmt.Sprintf(
		"Прости, пока рыбов не продаём, просто показываем, кросивое, да? %v",
		emoji.BeamingFaceWithSmilingEyes)
	<-ch
	return str
}

func removeInlineBlock(chatId int64, msgId int, bot *tgbotapi.BotAPI) {
	editedMsg := tgbotapi.NewEditMessageReplyMarkup(
		chatId,
		msgId,
		tgbotapi.InlineKeyboardMarkup{InlineKeyboard: make([][]tgbotapi.InlineKeyboardButton, 0, 0)})
	_, _ = bot.Send(editedMsg)
}
