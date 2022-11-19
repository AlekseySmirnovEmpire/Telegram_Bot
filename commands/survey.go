package commands

import (
	data2 "Telegram_Bot/data"
	"fmt"
	tgbotapi "github.com/Syfaro/telegram-bot-api"
	"github.com/enescakir/emoji"
	"strconv"
	"strings"
)

func survey(
	data *[]string,
	upd *tgbotapi.Update,
	bot *tgbotapi.BotAPI) (err error) {

	msg := tgbotapi.NewMessage(upd.CallbackQuery.Message.Chat.ID, "")
	switch (*data)[1] {
	case "New":
		msg.Text, err = newSurvey(data, &msg)
		break
	case "Pair":
		msg.Text = fmt.Sprintf(
			"Простите, пока не доступно %v, но мы уже страемся над тем, чтобы вы скоро смогли этим воспользоваться %v",
			emoji.ConfusedFace,
			emoji.WinkingFace)
		break
	case "Old":
		msg.Text = redactSurvey(data, &msg)
		break
	case "Restart":
		clearMessagesList(&(*data)[2], upd.CallbackQuery.Message.Chat.ID, bot)
		msg.Text, err = restartSurvey(data, upd, &msg)
	default:
		msg.Text, err = newQuest(data, &msg)
	}
	if err != nil {
		return err
	}

	err = editAndSendMessage(upd.CallbackQuery.Message.Chat.ID, bot, &(*data)[2], &msg)
	if err != nil {
		return err
	}

	return nil
}

func redactSurvey(data *[]string, msg *tgbotapi.MessageConfig) string {
	u, err := FindUser((*data)[2])
	if err != nil {
		return ""
	}

	if u.QuestCount < len(ql) {
		return "Вы ещё не прошли полностью анкету! Сперва завершите её!"
	}

	msg.ReplyMarkup = tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData(
				"Редактировать",
				fmt.Sprintf(
					"Pager:Init:%s:%s",
					(*data)[2],
					(*data)[3])),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData(
				"Пройти заново",
				fmt.Sprintf(
					"Survey:Restart:%s:%s",
					(*data)[2],
					(*data)[3])),
		))
	return "Нажмите \"Пройти заново\", если вы хотите снова пройти анкету. Если вы хотите изменить свой ответ на конкретный вопрос - нажмите \"Редактировать\"."
}

func restartSurvey(data *[]string, upd *tgbotapi.Update, msg *tgbotapi.MessageConfig) (str string, err error) {
	user, err := FindUser((*data)[2])
	if err != nil {
		return "", err
	}

	err = data2.DeleteAnswer(0, &(*data)[2], true)
	if err != nil {
		return "", err
	}

	user.QuestCount = 0
	str, err = sendQuestion(1, data, msg)
	if err != nil {
		return "", err
	}
	return str, nil
}

func newSurvey(data *[]string, msg *tgbotapi.MessageConfig) (string, error) {
	u, err := FindUser((*data)[2])
	if err != nil {
		return "", err
	}
	if u.QuestCount >= len(ql) {
		return "Вы уже проходили анкету!", nil
	} else if u.QuestCount > 0 {
		str, err1 := sendQuestion(u.QuestCount+1, data, msg)
		if err1 != nil {
			return "", err1
		}
		return str, nil
	}
	questID := u.QuestCount + 1

	msg.ReplyMarkup = tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData(
				"Начать",
				fmt.Sprintf(
					"Survey:%d:%s:%s",
					questID,
					(*data)[2],
					(*data)[3])),
		))
	return fmt.Sprintf(
		"%s, сейчас вы будете проходить тест, Вам будет предложено 3 варианта ответа на каждый вопрос: \"да\", \"нет\", \"возможно\". Если вы готовы, то нажмите \"Начать\"!",
		strings.Split(u.Name, " ")[0]), err
}

func newQuest(data *[]string, msg *tgbotapi.MessageConfig) (str string, err error) {
	questID, err := strconv.Atoi((*data)[1])
	if err != nil {
		str, err = continueQuest(data, msg)
		if err != nil {
			return "", err
		}
	}
	if str == "" {
		str, err = sendQuestion(questID, data, msg)
		if err != nil {
			return "", err
		}
	}

	return str, nil
}

func continueQuest(data *[]string, msg *tgbotapi.MessageConfig) (string, error) {
	strData := strings.Split((*data)[1], "_")
	questID, err := strconv.Atoi(strData[0])
	if err != nil {
		return "", err
	}
	var ans string
	switch strData[1] {
	case "yes":
		ans = "Да"
		break
	case "no":
		ans = "Нет"
		break
	default:
		ans = "Возможно"
		break
	}
	err = data2.InsertAnswer(&(*data)[2], &ans, questID)
	if err != nil {
		return "", err
	}

	u, err := FindUser((*data)[2])
	if err != nil {
		return "", err
	}
	u.QuestCount++

	// Анкета закончилась
	if questID+1 > len(ql) {
		msg.ReplyMarkup = tgbotapi.NewInlineKeyboardMarkup(
			tgbotapi.NewInlineKeyboardRow(
				tgbotapi.NewInlineKeyboardButtonData(
					"Выйти в меню",
					fmt.Sprintf(
						"Menu:quest:%s:%s",
						(*data)[2],
						(*data)[3])),
			))
		return "Спасибо, что прошли нашу анкету!", nil
	}

	str, err := sendQuestion(questID+1, data, msg)
	if err != nil {
		return "", err
	}

	return str, nil
}

func sendQuestion(questID int, data *[]string, msg *tgbotapi.MessageConfig) (string, error) {
	msg.ReplyMarkup = tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData(
				"Да",
				fmt.Sprintf(
					"Survey:%d_yes:%s:%s",
					questID,
					(*data)[2],
					(*data)[3])),
			tgbotapi.NewInlineKeyboardButtonData(
				"Возможно",
				fmt.Sprintf(
					"Survey:%d_maybe:%s:%s",
					questID,
					(*data)[2],
					(*data)[3])),
			tgbotapi.NewInlineKeyboardButtonData(
				"Нет",
				fmt.Sprintf(
					"Survey:%d_no:%s:%s",
					questID,
					(*data)[2],
					(*data)[3])),
		))

	return fmt.Sprintf("%d/%d: %s", questID, len(ql), ql[questID-1].Text), nil
}
