package commands

import (
	"Telegram_Bot/data"
	"Telegram_Bot/errors"
	tgbotapi "github.com/Syfaro/telegram-bot-api"
	"log"
	"reflect"
)

var users map[int64]*data.User

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

	for upd := range updates {
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

func FindUser(ID int64) (*data.User, error) {
	if usr, ok := users[ID]; ok {
		return usr, nil
	}
	return nil, errors.NotFound{Val: string(ID), Key: "users"}
}
