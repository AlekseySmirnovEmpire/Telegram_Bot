package commands

import (
	"Telegram_Bot/data"
	"Telegram_Bot/errors"
	tgbotapi "github.com/Syfaro/telegram-bot-api"
	"log"
	"reflect"
	"strings"
)

var (
	users   map[int64]*data.User
	ageConf map[int64]bool
)

// Listen listener.
func Listen(bot *tgbotapi.BotAPI) error {
	//Устанавливаем время обновления
	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	//Получаем обновления от бота
	updates, err := bot.GetUpdatesChan(u)
	if err != nil {
		return err
	}
	ageConf = make(map[int64]bool, 0)
	users = make(map[int64]*data.User, 0)

	for upd := range updates {
		if upd.CallbackQuery != nil {
			data := strings.Split(upd.CallbackQuery.Data, ":")
			if len(data) < 2 {
				continue
			}
			switch data[0] {
			case "Age_Confirm_Yes":
				err = ageConfirm(bot, &upd, &data[1], &data[0])
			case "Age_Confirm_No":
				err = ageConfirm(bot, &upd, &data[1], &data[0])
			default:
				_, _ = bot.Send(tgbotapi.NewMessage(upd.CallbackQuery.Message.Chat.ID, "Не распознал ответ!"))
			}

			if err != nil {
				_, _ = bot.Send(tgbotapi.NewMessage(upd.CallbackQuery.Message.Chat.ID, "Вы уже отвечали на этот вопрос!"))
			}
		}

		//There is no message.
		if upd.Message == nil {
			continue
		}

		//Message is text message (no video, sticker etc.).
		if reflect.TypeOf(upd.Message.Text).Kind() == reflect.String && upd.Message.Text != "" {
			switch upd.Message.Text {
			case "/start":
				err = start(bot, &upd)
			default:
				msg := tgbotapi.NewMessage(upd.Message.Chat.ID, upd.Message.Text)
				msg.Text = "Простите, но я не знаю такой команды!"
				_, _ = bot.Send(msg)
			}
		} else {
			msg := tgbotapi.NewMessage(upd.Message.Chat.ID, "Простите, не распознаю команду!")
			_, err = bot.Send(msg)
		}

		if err != nil {
			log.Printf("Error: %s", err.Error())
		}
	}

	return nil
}

func AddUser(usr *data.User, ID int64) {
	if len(users) == 0 {
		users = make(map[int64]*data.User, 0)
	}
	users[ID] = usr
}

func FindUserInArray(ID int64) (*data.User, error) {
	if len(users) == 0 {
		return nil, errors.NotFound{Val: string(ID), Key: "users"}
	}
	if usr, ok := users[ID]; ok {
		return usr, nil
	}
	return nil, errors.NotFound{Val: string(ID), Key: "users"}
}

func CheckAge(ID int64) bool {
	if ac, ok := ageConf[ID]; ok {
		return ac
	}

	return false
}
