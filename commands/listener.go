package commands

import (
	"Telegram_Bot/data"
	"Telegram_Bot/errors"
	"fmt"
	tgbotapi "github.com/Syfaro/telegram-bot-api"
	"github.com/enescakir/emoji"
	"log"
	"reflect"
	"strconv"
)

type users map[string]*data.User

var (
	um           users
	ageConfMesID map[string]int
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
	ageConfMesID = make(map[string]int, 0)

	for upd := range updates {
		if upd.Message != nil {
			msg := tgbotapi.NewMessage(upd.Message.Chat.ID, upd.Message.Text)

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

			// Если получили ошибку, то отвечаем что что-то не так.
			if err != nil {
				log.Printf("Error: %s", err.Error())
				msg = tgbotapi.NewMessage(upd.Message.Chat.ID, upd.Message.Text)
				msg.Text = "Простите, что-то я сегодня туплю!"
				_, _ = bot.Send(msg)
				continue
			}

			// Отправляем полученное сообщение
			_, _ = bot.Send(msg)

			// Если юзер не подтверждал возраст
			if IsUserAuth(strconv.Itoa(upd.Message.From.ID)) && !um[strconv.Itoa(upd.Message.From.ID)].AgeConfirmed {
				msg = tgbotapi.NewMessage(upd.Message.Chat.ID, upd.Message.Text)
				msg.Text = AgeConfirm(&msg, &upd)
				msg2, err2 := bot.Send(msg)
				if err2 == nil {
					userID := strconv.Itoa(upd.Message.From.ID)
					if _, ok := ageConfMesID[userID]; ok {
						delete(ageConfMesID, userID)
					}
					ageConfMesID[userID] = msg2.MessageID
				}
			}
		} else {
			// Мы получили ответ через кнопку
			if upd.CallbackQuery != nil {
				err = callBack(&upd, bot)
				if err != nil {
					_, _ = bot.Send(tgbotapi.NewMessage(
						upd.CallbackQuery.Message.Chat.ID,
						fmt.Sprintf("Упс, что-то пошло не так %v!", emoji.ThinkingFace)))
				}
			}
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
