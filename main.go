package main

import (
	bot2 "Telegram_Bot/bot"
	"Telegram_Bot/commands"
	"Telegram_Bot/config"
	"Telegram_Bot/db"
	"github.com/joho/godotenv"
	"log"
)

func init() {
	if err := godotenv.Load(); err != nil {
		log.Print("No .env file found")
	}
}

func main() {
	//init configs
	err := config.Init()
	if err != nil {
		log.Println(err.Error())
		return
	}

	//init db
	defer db.CloseDB()
	err = db.InitDB()
	if err != nil {
		log.Println(err.Error())
		return
	}

	//Init bot itself
	bot, err := bot2.CreateBot()
	if err != nil {
		log.Println(err.Error())
		return
	}

	log.Println("Bot is running ....")

	err = commands.Listen(bot)
	if err != nil {
		log.Printf("Bot SHIT DOWN! Error: %s", err.Error())
	}
}
