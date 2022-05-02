package simplectx

import bot "gopkg.in/telebot.v3"

type Sticker string

func (s Sticker) Send(b *bot.Bot, r bot.Recipient, o *bot.SendOptions) (*bot.Message, error) {
	return b.Send(r, &bot.Sticker{File: bot.File{FileID: string(s)}}, o)
}
