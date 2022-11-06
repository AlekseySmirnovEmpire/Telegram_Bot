package config

import (
	"Telegram_Bot/myErrors"
	"log"
	"os"
)

var (
	Token  string
	IsDev  bool
	ConStr string
)

func Init() error {
	log.Println("Reading \".env\" file ....")

	token, exist := os.LookupEnv("Telegram_Token")
	if !exist || token == "" {
		return myErrors.EmptyFile{Val: ".env"}
	}
	conStr, exist := os.LookupEnv("Connection_String")
	if !exist || conStr == "" {
		return myErrors.EmptyFile{Val: ".env"}
	}
	env, exist := os.LookupEnv("Environment")
	if !exist || env == "" {
		return myErrors.EmptyFile{Val: ".env"}
	}
	Token = token
	ConStr = conStr
	if env == "Production" {
		IsDev = false
	} else {
		IsDev = true
	}

	log.Println("Complete reading \".env\" file success!")

	return nil
}
