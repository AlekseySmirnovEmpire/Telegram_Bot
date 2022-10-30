package commands

import (
	tgbotapi "github.com/Syfaro/telegram-bot-api"
	"log"
	"reflect"
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

	for upd := range updates {
		//There is no message.
		if upd.Message == nil {
			continue
		}

		//Message is text message (no video, sticker etc.).
		if reflect.TypeOf(upd.Message.Text).Kind() == reflect.String && upd.Message.Text != "" {
			switch upd.Message.Text {
			case "/start":
				err = start(bot, upd)
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
