package bot

import (
	"Telegram_Bot/config"
	"Telegram_Bot/errors"
	tgbotapi "github.com/Syfaro/telegram-bot-api"
	"log"
)

// CreateBot create new telegram bot.
func CreateBot() (*tgbotapi.BotAPI, error) {
	log.Println("Start connecting to Telegram with token ....")

	bot, err := tgbotapi.NewBotAPI(config.Token)
	if err != nil {
		return nil, errors.NoConnection{Val: "Telegram", Err: err.Error()}
	}

	if config.IsDev {
		bot.Debug = true
		log.Println("Bot debug set TRUE!")
	}

	log.Printf("Authorized on account %s!", bot.Self.UserName)
	log.Println("Bot created success!")

	return bot, nil
}
