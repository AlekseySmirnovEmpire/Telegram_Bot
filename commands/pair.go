package commands

import (
	"Telegram_Bot/data"
	"Telegram_Bot/myErrors"
	"fmt"
	"log"

	tgbotapi "github.com/Syfaro/telegram-bot-api"
	"github.com/enescakir/emoji"
)

func pair(data *[]string,
	upd *tgbotapi.Update,
	bot *tgbotapi.BotAPI) (err error) {

	msg := tgbotapi.NewMessage(upd.CallbackQuery.Message.Chat.ID, "")

	switch (*data)[1] {
	case "Init":
		clearMessagesList(&(*data)[2], upd.CallbackQuery.Message.Chat.ID, bot)
		msg.Text = initPair(&msg, data)
		break
	case "New":
		clearMessagesList(&(*data)[2], upd.CallbackQuery.Message.Chat.ID, bot)
		msg.Text, err = newPairSurvey(data, &msg)

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

func initPair(msg *tgbotapi.MessageConfig, data *[]string) (str string) {
	msg.ReplyMarkup = tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData(
				"Создать новую",
				fmt.Sprintf(
					"Pair:New:%s:%s",
					(*data)[2],
					(*data)[3])),
			tgbotapi.NewInlineKeyboardButtonData(
				"Указать пару",
				fmt.Sprintf(
					"Pair:Mark:%s:%s",
					(*data)[2],
					(*data)[3])),
		))

	str = "Укажите свою пару, если она у вас уже есть, или создайте новую!"
	return
}

func newPairSurvey(dataFromQuery *[]string, msg *tgbotapi.MessageConfig) (str string, err error) {
	if user, ok := um[(*dataFromQuery)[2]]; ok {
		// Юзер ещё не прошёл свою анкету (+)
		if user.QuestCount != 0 && user.QuestCount < len(ql) {
			*dataFromQuery = append(*dataFromQuery, "NeedID")
			str, err = newSurvey(dataFromQuery, msg)

			str = "Для начала пройдите свою анкету, давайте продолжим с того места, где вы остановились:\n" + str

			return
		}

		// Юзер уже проходил свою анкету и надо вернуть ему её айди.
		if user.QuestCount == len(ql) {
			id, err1 := data.ReturnUserQuestinaryId(&(*dataFromQuery)[2], true)
			if err1 != nil {
				str = "Извините, ваши ответы сохранились, но сейчас ваша анкета не доступна! Обратитесь в тех поддержку!"
				log.Printf("Error with returning user anket ID, err: %s", err1.Error())
			} else {
				str = fmt.Sprintf("ID вашей анкеты: `%s`", id)
			}

			msg.ReplyMarkup = tgbotapi.NewInlineKeyboardMarkup(
				tgbotapi.NewInlineKeyboardRow(
					tgbotapi.NewInlineKeyboardButtonData(
						"Выйти в меню",
						fmt.Sprintf(
							"Menu:quest:%s:%s",
							(*dataFromQuery)[2],
							(*dataFromQuery)[3])),
				))

			return
		}

		str = fmt.Sprintf("Укажите ID анкеты вашей пары, пришлите готовый или попросите свою пару вам отправить его %s:", emoji.WinkingFace)
		waitForId[(*dataFromQuery)[2]] = true

		return
	} else {
		err = myErrors.NotFound{Val: (*dataFromQuery)[2], Key: "User list", Err: err}
		return
	}
}
