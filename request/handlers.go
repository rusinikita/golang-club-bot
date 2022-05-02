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
	AnnounceLink     string `env:"ANNOUNCE_LINK,notEmpty"`
	QnALink          string `env:"QNA_LINK,notEmpty"`
	AdminID          int64  `env:"ADMIN_ID,required"`
	GroupChatID      int64  `env:"GROUP_CHAT_ID,required"`
	SendNotification bool   `env:"SEND_NOTIFICATION"`
}

func newConfig() config {
	cfg := config{}
	if err := env.Parse(&cfg); err != nil {
		log.Fatalf("%+v", err)
	}

	return cfg
}

func Setup(db airtable.Airtable, b *bot.Bot) {
	cfg := newConfig()

	goBtns := &bot.ReplyMarkup{}
	goButton := goBtns.Data("Go 💪", "go")
	goBtns.Inline(goBtns.Row(goButton))

	b.Handle("/start", func(c bot.Context) error {
		return simplectx.Wrap(c, func(c bot.Context, sc *simplectx.Context) {
			ctx := context.Background()

			sc.Send(hello)
			sc.Send(fmt.Sprintf(hiText, cfg.AnnounceLink), &bot.SendOptions{
				DisableWebPagePreview: true,
				ReplyMarkup:           goBtns,
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
			if c.Chat().ID == cfg.GroupChatID {
				return
			}

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

	// notification sending
	if !cfg.SendNotification {
		return
	}

	var requests []Request
	ctx := context.Background()
	e(db.List(ctx, &requests, airtable.View("NotSend")))

	for _, r := range requests {
		switch r.Status {
		case None:
			send(b, r, sad)
			send(b, r, fmt.Sprintf(remind, 11), goBtns)

		case Declined:
			send(b, r, sorry)
			send(b, r, decline)
			send(b, r, r.DeclineMessage)

		case Accepted:
			link, err := b.CreateInviteLink(bot.ChatID(cfg.GroupChatID), &bot.ChatInviteLink{
				Name:        r.Name + " link",
				MemberLimit: 1,
			})
			e(err)

			btns := &bot.ReplyMarkup{}
			btn := btns.URL("👌 let's go", link.InviteLink)
			btns.Inline(btns.Row(btn))

			send(b, r, welcome)
			send(b, r, accept, btns)
		}

		e(db.Patch(ctx, Request{
			RecordID: r.RecordID,
			Send:     true,
		}))

		fmt.Println("done", r.UserID, r.Name)

		time.Sleep(500 * time.Millisecond)
	}

	fmt.Println("all done")
}

func send(b *bot.Bot, to bot.Recipient, what interface{}, opts ...interface{}) {
	_, err := b.Send(to, what, opts...)
	e(err)
}

func e(err error) {
	if err == nil {
		return
	}

	log.Fatal(err)
}
