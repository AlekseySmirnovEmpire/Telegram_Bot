package commands

import (
	"Telegram_Bot/data"
	"Telegram_Bot/myErrors"
	"errors"
	"fmt"
	"log"
	"reflect"
	"strconv"

	tgbotapi "github.com/Syfaro/telegram-bot-api"
	"github.com/enescakir/emoji"
)

type users map[string]*data.User
type questions []*data.Question
type userPairs map[string][]*data.User
type waitingForUsersForID map[string]bool

var (
	um              users
	messageToDelete map[string]int
	ql              questions
	pagerMes        map[string]int
	pairs           userPairs
	waitForId       waitingForUsersForID
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
	pagerMes = make(map[string]int, 0)

	//ПРогружаем всех пар.
	log.Println("Loading pairs ....")
	for _, user := range um {
		ul, err1 := data.LoadPairs(user.ID)
		if err1 != nil {
			log.Printf("There is error while trying to get pairs for user %s", user.Key)
			continue
		}

		if ul != nil && len(ul) > 0 {
			pairs[user.Key] = ul
		}
	}
	log.Printf("Loaded %d pairs!", len(pairs))

	waitForId = make(map[string]bool, 0)

	ql, err = data.InitQuestions()
	if err != nil || ql == nil {
		return myErrors.NotFound{Val: "questions", Key: err.Error()}
	}

	for upd := range updates {
		if upd.Message != nil {
			userID := strconv.Itoa(upd.Message.From.ID)

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
		} else if upd.InlineQuery != nil {
			// Наш инлайн затригерился
			var articles []interface{}
			uId := strconv.Itoa(upd.InlineQuery.From.ID)

			if _, ok := um[uId]; ok {
				mes := fmt.Sprintf(
					"Эй, привет!\nЯ прошёл опрос в анкете и приглашаю тебя пройти парный опрос %s\nЗаходи в бота, жми по кнопке ниже %s",
					emoji.WinkingFace, emoji.BackhandIndexPointingDown)
				msg := tgbotapi.NewInlineQueryResultArticleMarkdown(
					"Pair",
					"Парный опрос",
					mes)
				keyboard := tgbotapi.NewInlineKeyboardMarkup(
					tgbotapi.NewInlineKeyboardRow(
						tgbotapi.NewInlineKeyboardButtonURL(
							"Перейти в бота",
							fmt.Sprintf("https://t.me/%s", bot.Self.UserName))))
				msg.ReplyMarkup = &keyboard
				articles = append(articles, msg)
			} else {
				msg := tgbotapi.NewInlineQueryResultArticleMarkdown(
					upd.InlineQuery.ID,
					"Ошибка!",
					"Что-то пошло не так! Проверьте, подтвердили ли вы свой возраст в боте!")
				articles = append(articles, msg)
			}

			inlineConfig := tgbotapi.InlineConfig{
				InlineQueryID: upd.InlineQuery.ID,
				IsPersonal:    true,
				CacheTime:     0,
				Results:       articles,
			}

			_, err2 := bot.AnswerInlineQuery(inlineConfig)
			if err2 != nil {
				log.Println(err)
			}
			chatId, err2 := strconv.ParseInt(uId, 10, 64)
			if err2 == nil {
				message := tgbotapi.NewMessage(chatId, fmt.Sprintf("Парный опрос отправлен %s", emoji.ThumbsUp))
				err2 = editAndSendMessage(chatId, bot, &uId, &message)
				if err2 != nil {
					log.Println(err2)
				}
			}
		} else {
			// Мы получили ответ через кнопку
			if upd.CallbackQuery != nil {

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

func editAndSendMessage(chatID int64, bot *tgbotapi.BotAPI, key *string, msg *tgbotapi.MessageConfig) error {
	if msgID, ok := messageToDelete[*key]; ok {
		editMes := tgbotapi.NewEditMessageText(chatID, msgID, msg.Text)
		ms, err := bot.Send(editMes)
		if err != nil {
			return err
		}
		delete(messageToDelete, *key)
		if msg.ReplyMarkup != nil {
			editMes1 := tgbotapi.NewEditMessageReplyMarkup(chatID, msgID, msg.ReplyMarkup.(tgbotapi.InlineKeyboardMarkup))
			ms, err = bot.Send(editMes1)
			if err != nil {
				return err
			}
			messageToDelete[*key] = ms.MessageID
		}
	} else {
		ms, err := bot.Send(msg)
		if err != nil {
			return err
		}
		if msg.ReplyMarkup != nil {
			messageToDelete[*key] = ms.MessageID
		}
	}

	return nil
}

func clearMessagesList(key *string, chatID int64, bot *tgbotapi.BotAPI) {
	if msgID, ok := messageToDelete[*key]; ok {
		removeInlineBlock(chatID, msgID, bot)
		delete(messageToDelete, *key)
	}
}
