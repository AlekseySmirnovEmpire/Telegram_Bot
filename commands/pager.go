package commands

import (
	data2 "Telegram_Bot/data"
	"Telegram_Bot/myErrors"
	"fmt"
	tgbotapi "github.com/Syfaro/telegram-bot-api"
	"strconv"
	"strings"
)

func pager(data *[]string,
	upd *tgbotapi.Update,
	bot *tgbotapi.BotAPI) (err error) {
	msg := tgbotapi.NewMessage(upd.CallbackQuery.Message.Chat.ID, "")

	switch (*data)[1] {
	case "Init":
		clearMessagesList(&(*data)[2], upd.CallbackQuery.Message.Chat.ID, bot)
		msg.Text = initPager(data, 0, &msg)
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

func recognize(
	data *[]string,
	msg *tgbotapi.MessageConfig,
	bot *tgbotapi.BotAPI,
	upd *tgbotapi.Update) (string, error) {

	req := strings.Split((*data)[1], "_")
	if len(req) < 2 || len(req) > 3 {
		return "", myErrors.NotConfirmed{Key: "request", Val: "pager"}
	}

	chatID, err := strconv.ParseInt((*data)[3], 10, 64)
	if err != nil {
		return "", myErrors.NotConfirmed{Key: "request", Val: "pager"}
	}

	var str string
	removeInlineBlock(chatID, pagerMes[(*data)[2]], bot)
	delete(pagerMes, (*data)[2])
	questID, err := strconv.Atoi(req[1])
	if err != nil {
		return "", err
	}

	switch req[0] {
	case "Next":
		str = initPager(data, questID, msg)
		break
	case "Prev":
		str = initPager(data, questID, msg)
		break
	case "Quest":
		str = refactorQuest(data, questID, msg)
		break
	case "Refactor":
		str, err = refactorAnswer(data, questID, &req[2])
		err = initMainMenu(upd, bot, upd.CallbackQuery.Message.Chat.ID)
		break
	}

	if err != nil {
		return "", err
	}

	return str, nil
}

func refactorAnswer(data *[]string, questID int, ans *string) (string, error) {
	var a string
	switch *ans {
	case "yes":
		a = "Да"
		break
	case "no":
		a = "Нет"
		break
	default:
		a = "Возможно"
	}

	err := data2.UpdateQuestion(&(*data)[2], &a, questID+1)
	if err != nil {
		return "", err
	}

	return fmt.Sprintf("Вы успешно заминили ответ на \"%s\"", a), nil
}

func refactorQuest(data *[]string, questID int, msg *tgbotapi.MessageConfig) string {
	msg.ReplyMarkup = tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData(
				"Да",
				fmt.Sprintf(
					"Pager:Refactor_%d_yes:%s:%s",
					questID,
					(*data)[2],
					(*data)[3])),
			tgbotapi.NewInlineKeyboardButtonData(
				"Возможно",
				fmt.Sprintf(
					"Pager:Refactor_%d_maybe:%s:%s",
					questID,
					(*data)[2],
					(*data)[3])),
			tgbotapi.NewInlineKeyboardButtonData(
				"Нет",
				fmt.Sprintf(
					"Pager:Refactor_%d_no:%s:%s",
					questID,
					(*data)[2],
					(*data)[3])),
		))

	return fmt.Sprintf("%d/%d: %s", questID+1, len(ql), ql[questID].Text)
}

func initPager(data *[]string,
	start int,
	msg *tgbotapi.MessageConfig) string {

	str := "Выберите номер вопроса из списка для редактирования вашего ответа:\n"

	var rows []tgbotapi.InlineKeyboardButton
	if start-5 >= 0 {
		rows = append(rows, tgbotapi.NewInlineKeyboardButtonData(
			"<<",
			fmt.Sprintf(
				"Pager:Prev_%d:%s:%s",
				start-5,
				(*data)[2],
				(*data)[3])))
	}

	ind := start
	for ; ind < start+5 && ind <= len(ql)-1; ind++ {
		str += fmt.Sprintf("%d. %s;\n", ind+1, ql[ind].Text)
		rows = append(rows, tgbotapi.NewInlineKeyboardButtonData(
			fmt.Sprintf("%d", ind+1),
			fmt.Sprintf(
				"Pager:Quest_%d:%s:%s",
				ql[ind].ID-1,
				(*data)[2],
				(*data)[3])),
		)
	}

	if ind < len(ql)-1 {
		rows = append(rows, tgbotapi.NewInlineKeyboardButtonData(
			">>",
			fmt.Sprintf(
				"Pager:Next_%d:%s:%s",
				ind,
				(*data)[2],
				(*data)[3])))
	}

	msg.ReplyMarkup = tgbotapi.NewInlineKeyboardMarkup(rows)
	return str
}
