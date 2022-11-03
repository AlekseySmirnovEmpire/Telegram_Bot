package commands

import (
	"Telegram_Bot/data"
	"Telegram_Bot/errors"
	tgbotapi "github.com/Syfaro/telegram-bot-api"
	"log"
	"reflect"
	"strconv"
)

type users map[string]*data.User

var (
	um users
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

	//Прогружаем всех пользователей из БД
	um, err = data.InitUsers()
	if err != nil {
		return err
	}

	for upd := range updates {
		msg := tgbotapi.NewMessage(upd.Message.Chat.ID, upd.Message.Text)
		//if upd.CallbackQuery != nil {
		//	data := strings.Split(upd.CallbackQuery.Data, ":")
		//	if len(data) < 2 {
		//		continue
		//	}
		//	switch data[0] {
		//	case "Age_Confirm_Yes":
		//		err = ageConfirm(bot, &upd, &data[1], &data[0])
		//	case "Age_Confirm_No":
		//		err = ageConfirm(bot, &upd, &data[1], &data[0])
		//	default:
		//		_, _ = bot.Send(tgbotapi.NewMessage(upd.CallbackQuery.Message.Chat.ID, "Не распознал ответ!"))
		//	}
		//
		//	if err != nil {
		//		_, _ = bot.Send(tgbotapi.NewMessage(upd.CallbackQuery.Message.Chat.ID, "Вы уже отвечали на этот вопрос!"))
		//	}
		//}

		//There is no message.
		if upd.Message == nil {
			continue
		}

		// Message is text message (no video, sticker etc.).
		if reflect.TypeOf(upd.Message.Text).Kind() == reflect.String && upd.Message.Text != "" {
			switch upd.Message.Text {
			case "/start":
				msg.Text, err = start(&upd)
				break
			default:
				msg.Text = getRandomShit(upd.Message)
			}
		} else {
			msg.Text = getRandomShit(upd.Message)
		}

		if err != nil {
			log.Printf("Error: %s", err.Error())
			msg = tgbotapi.NewMessage(upd.Message.Chat.ID, upd.Message.Text)
			msg.Text = "Простите, что-то я сегодня туплю!"
			_, _ = bot.Send(msg)
			continue
		}

		_, _ = bot.Send(msg)
		if !um[strconv.Itoa(upd.Message.From.ID)].AgeConfirmed {
			msg = tgbotapi.NewMessage(upd.Message.Chat.ID, upd.Message.Text)
			msg.Text = AgeConfirm(&msg, &upd)
			_, _ = bot.Send(msg)
		}
	}

	return nil
}

func IsUserAuth(userID string) bool {
	_, ok := um[userID]
	return ok
}

func FindUser(userID string) (*data.User, error) {
	if u, ok := um[userID]; ok {
		return u, nil
	}

	return nil, errors.NotFound{Val: "users map", Key: "no user in map"}
}

func addUser(u *data.User) {
	um[u.Key] = u
}
