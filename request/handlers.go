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
	AdminID      int64  `env:"ADMIN_ID,required"`
	GroupChatID  int64  `env:"GROUP_CHAT_ID,required"`
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

	b.Handle("/notify", func(c bot.Context) error {
		return simplectx.Wrap(c, func(c bot.Context, sc *simplectx.Context) {
			if c.Sender().ID != cfg.AdminID {
				sc.Reply("permission denied")

				return
			}

			var requests []Request
			filter := airtable.Filter(Request{
				Send: false,
			})
			ctx := context.Background()
			sc.Error(db.List(ctx, &requests, filter))

			acceptedCount := 0
			for _, r := range requests {
				if r.Status == Accepted {
					acceptedCount++
				}
			}

			for _, r := range requests {
				switch r.Status {
				case None:
					sc.SendTo(r, sad)
					sc.SendTo(r, fmt.Sprintf(remind, acceptedCount), goButton)

				case Declined:
					sc.SendTo(r, sorry)
					sc.SendTo(r, decline)
					sc.SendTo(r, r.DeclineMessage)

				case Accepted:
					link, err := c.Bot().CreateInviteLink(bot.ChatID(cfg.GroupChatID), &bot.ChatInviteLink{
						Name:        r.Name + " link",
						MemberLimit: 1,
					})
					sc.Error(err)

					btns := &bot.ReplyMarkup{}
					btn := btns.URL("Согласен, let's go", link.InviteLink)
					btns.Inline(btns.Row(btn))

					sc.SendTo(r, welcome)
					sc.SendTo(r, accept, btns)
				}

				sc.Error(db.Patch(ctx, Request{
					RecordID: r.RecordID,
					Send:     true,
				}))

				fmt.Println("done", r.UserID, r.Name)

				time.Sleep(500 * time.Millisecond)
			}
		})
	})
}
