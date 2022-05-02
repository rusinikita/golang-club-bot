package request

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/caarlos0/env/v6"
	bot "gopkg.in/telebot.v3"

	"github.com/rusinikita/gogoClub/airtable"
	"github.com/rusinikita/gogoClub/simplectx"
)

type config struct {
	AnnounceLink string `env:"ANNOUNCE_LINK,notEmpty"`
	QnALink      string `env:"QNA_LINK,notEmpty"`
}

func Setup(db airtable.Airtable, b *bot.Bot) {
	cfg := config{}
	if err := env.Parse(&cfg); err != nil {
		log.Fatalf("%+v", err)
	}

	btns := &bot.ReplyMarkup{}
	goButton := btns.Data("Go 💪", "go")
	btns.Inline(btns.Row(goButton))

	b.Handle("/start", func(c bot.Context) error {
		return simplectx.Wrap(c, func(c bot.Context, sc *simplectx.Context) {
			ctx := context.Background()

			sc.Send(hello)
			sc.Send(fmt.Sprintf(hiText, cfg.AnnounceLink), &bot.SendOptions{
				DisableWebPagePreview: true,
				ReplyMarkup:           btns,
			})

			var users []Request
			err := db.List(ctx, &users, airtable.Filter(NewID(c.Sender())))
			if err != nil {
				log.Println(err)
			}

			if len(users) > 0 {
				return
			}

			sc.Error(db.Create(ctx, New(c.Sender())))
		})
	})

	b.Handle(&goButton, func(c bot.Context) error {
		return simplectx.Wrap(c, func(c bot.Context, sc *simplectx.Context) {
			interval := 2 * time.Second
			sc.Send(step1, bot.NoPreview)

			time.Sleep(interval)
			sc.Send(step2, bot.NoPreview, bot.Silent)

			time.Sleep(interval)
			sc.Send(step3, bot.NoPreview, bot.Silent)

			time.Sleep(interval)
			sc.Send(step4, bot.NoPreview, bot.Silent)

			time.Sleep(interval)
			sc.Send(step5, bot.NoPreview, bot.Silent)

			time.Sleep(interval)
			sc.Send(fmt.Sprintf(step6, cfg.QnALink), bot.NoPreview, bot.Silent)
		})
	})

	b.Handle(bot.OnText, func(c bot.Context) error {
		return simplectx.Wrap(c, func(c bot.Context, sc *simplectx.Context) {
			link := ""
			for _, entity := range c.Message().Entities {
				if entity.Type != bot.EntityURL && entity.Type != bot.EntityTextLink {
					continue
				}

				link = entity.URL
				if link == "" {
					link = c.Message().EntityText(entity)
				}

				break
			}

			if link == "" {
				sc.Send(wtf)
				sc.Send(fmt.Sprintf(noLinkText, cfg.QnALink), bot.NoPreview)

				return
			}

			var (
				ctx   = context.Background()
				users []Request
				user  = New(c.Sender())
			)

			sc.Error(db.List(ctx, &users, airtable.Filter(NewID(c.Sender()))))

			if len(users) != 1 {
				user := New(c.Sender())
				user.RepositoryLink = link
				sc.Error(db.Create(context.Background(), user))
			} else {
				user = users[0]
				user.RepositoryLink = link
				sc.Error(db.Patch(context.Background(), user))
			}

			sc.Send(wellDone)
			sc.Send(doneText)
		})
	})
}

const (
	hello    simplectx.Sticker = "CAACAgQAAxkBAAOpYmmXKS3ykycZk2qrR97R2_jTLKwAAswAA845CA3fZ3xlfkS5ZCQE"   // 🙃
	wtf      simplectx.Sticker = "CAACAgQAAxkBAAOrYmmXgeYtrZiN2IuUckR854EheykAApkAA845CA0jIAABUzXpH78kBA" // 🤨
	wellDone simplectx.Sticker = "CAACAgQAAxkBAAOtYmmX5-cveqGjl44BirOjkuy1cz4AApcAA845CA1AIS58gGBWGiQE"   // 👍
)

const (
	hiText = `Привет!

Если ты хочешь поучаствовать <a href="%s">в клубе изучения Go</a>, то ты по адресу.

Пройди простое задание, чтобы доказать мотивацию, и сможешь учиться GoLang с теми, кто поможет справиться с трудностями и подскажет что делать после.`
	step1 = `1. Создай профиль на github.com

Заполни настоящее имя, фото и город. Напиши в bio текущее место учебы/работы (направление и курс тоже)`
	step2 = `2. <a href="https://drive.google.com/file/d/1-8AQtU5WuftQrUioXYkp0bY2K20t9vM3/view?usp=sharing">Создай репозиторий</a>

Название kgl-go-learing (с дополнениями если занято)
README файл (нужно поставить галочку)

Этот репозиторий станет твоим портфолио, тетрадкой с заданиями и поможет систематизировать знания`
	step3 = `3. Напиши в README файле

1. О себе
2. Цель изучения. Почему ты хочешь научиться go или программированию ("почему бы и нет" - валидный ответ, но нужно подробнее расписать)
3. Почему уверен, что не бросишь занятия и не потратишь время менторов впустую
4. Ожидания. Как скоро и что ты хочешь получить от участия.
5. Перечисли вопросы, которые хотел бы обсудить.

Используй # для создания заголовков`
	step4 = `4. Просто отправь мне ссылку на получившийся репозиторий в сообщении.

Например, https://github.com/rusinikita/mindful-bot`
	step5      = `5. Дождись конца проверки и ответа`
	step6      = `Это всё, жду ссылку. Если какие-то проблемы, предложения, что-то не нравится - напиши <a href="%s">в этом чате</a>`
	noLinkText = `Ты не прислал ссылку. Если хочешь обсудить что-то - напиши <a href="%s">в этом чате</a>`
	doneText   = `Супер. Как только проверят задание, я отправлю тебе результат.

Если хочешь поменять ссылку, отправь новую.`
)