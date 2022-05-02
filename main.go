package main

import (
	"fmt"
	"log"

	"github.com/caarlos0/env/v6"
	"github.com/joho/godotenv"
	bot "gopkg.in/telebot.v3"

	"github.com/rusinikita/gogoClub/airtable"
	"github.com/rusinikita/gogoClub/request"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Println("Error loading .env file")
	}

	db := airtable.New()
	b := newBot()

	request.Setup(db, b)

	b.Start()
}

func newBot() (b *bot.Bot) {
	cfg := botConfig{}
	if err := env.Parse(&cfg); err != nil {
		log.Fatalf("%+v", err)
	}

	prefs := bot.Settings{
		Token: cfg.Token,
		OnError: func(err error, c bot.Context) {
			log.Println(err)

			c.Send(fmt.Sprintf(
				"Пожалуйста, сообщите об ошибке <a href=\"%s\">в этом чате</a>\n\n%s",
				cfg.QnALink,
				err.Error(),
			), bot.NoPreview)
		},
		ParseMode: bot.ModeHTML,
	}

	b, err := bot.NewBot(prefs)
	if err != nil {
		log.Fatalf("%+v", err)
	}

	return b
}

type botConfig struct {
	Token   string `env:"BOT_TOKEN,notEmpty"`
	QnALink string `env:"QNA_LINK,notEmpty"`
}
