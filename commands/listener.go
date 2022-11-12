package commands

import (
	"Telegram_Bot/data"
	"Telegram_Bot/myErrors"
	"errors"
	"fmt"
	tgbotapi "github.com/Syfaro/telegram-bot-api"
	"github.com/enescakir/emoji"
	"log"
	"reflect"
	"strconv"
)

type users map[string]*data.User
type questions []*data.Question

var (
	um              users
	messageToDelete map[string]int
	ql              questions
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
	messageToDelete = make(map[string]int, 0)

	ql, err = data.InitQuestions()
	if err != nil || ql == nil {
		return myErrors.NotFound{Val: "questions", Key: err.Error()}
	}

	for upd := range updates {
		if upd.Message != nil {
			userID := strconv.Itoa(upd.Message.From.ID)

			// удаляем кнопки у прежнего сообщения
			clearMessagesList(&userID, upd.Message.Chat.ID, bot)

			msg := tgbotapi.NewMessage(upd.Message.Chat.ID, "")
			isInit := false

			// Message is text message (no video, sticker etc.).
			if reflect.TypeOf(upd.Message.Text).Kind() == reflect.String && upd.Message.Text != "" {
				switch upd.Message.Text {
				case "/start":
					msg.Text, err = start(&upd)
					break
				case "/menu":
					err = initMainMenu(&upd, bot, upd.Message.Chat.ID)
					isInit = true
					var noConf myErrors.NotConfirmed
					var notFound myErrors.NotFound
					if errors.As(err, &noConf) {
						msg.Text = noConf.Error()
						err = nil
					} else if errors.As(err, &notFound) {
						msg.Text = "Вы не прошли проверку на возраст!"
						err = nil
					}
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
			if msg.Text != "" {
				_, _ = bot.Send(msg)
			} else if !isInit {
				_ = initMainMenu(&upd, bot, upd.Message.Chat.ID)
			}

			// Если юзер не подтверждал возраст
			if IsUserAuth(strconv.Itoa(upd.Message.From.ID)) && !um[strconv.Itoa(upd.Message.From.ID)].AgeConfirmed {
				msg = tgbotapi.NewMessage(upd.Message.Chat.ID, upd.Message.Text)
				msg.Text = AgeConfirm(&msg, &upd)
				msg2, err2 := bot.Send(msg)
				if err2 == nil {
					messageToDelete[userID] = msg2.MessageID
				}
			}
		} else {
			// Мы получили ответ через кнопку
			if upd.CallbackQuery != nil {
				userID := strconv.Itoa(upd.CallbackQuery.From.ID)

				// удаляем скнопки из предыдущего сообщения
				clearMessagesList(&userID, upd.CallbackQuery.Message.Chat.ID, bot)

				// слушаем ответ от кнопки
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

	return nil, myErrors.NotFound{Val: "users map", Key: "no user in map", Err: errors.New("no users")}
}

func addUser(u *data.User) {
	um[u.Key] = u
}

func clearMessagesList(key *string, chatID int64, bot *tgbotapi.BotAPI) {
	if msgID, ok := messageToDelete[*key]; ok {
		removeInlineBlock(chatID, msgID, bot)
		delete(messageToDelete, *key)
	}
}
