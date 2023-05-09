package main

import (
	"log"
	"os"

	"github.com/muskelo/translator_bot/internal/bot"
)

func main() {
	settings := bot.Settings{
		DatabaseConnStr:   os.Getenv("DB_CONNSTR"),
		LibreTranslateUrl: os.Getenv("LIBRETRANSLATE_URL"),
		BotToken:          os.Getenv("BOT_TOKEN"),
		Defaults: bot.Defaults{
			PrimaryLanguage:   os.Getenv("DEFAULT_PRIMARY_LANG"),
			SecondaryLanguage: os.Getenv("DEFAULT_SECONDARY_LANG"),
		},
	}

	b, err := bot.New(settings)
	if err != nil {
		log.Fatalf("Can't create bot: %v\n", err)
	}
	b.Init()
	b.Start()
}
