package simplectx

import (
	bot "gopkg.in/telebot.v3"
)

type Context struct {
	c    bot.Context
	done chan error
}

func Wrap(c bot.Context, f func(c bot.Context, sc *Context)) error {
	sc := &Context{c: c, done: make(chan error)}

	go func() {
		f(c, sc)
		sc.done <- nil
	}()

	return <-sc.done
}

func (c *Context) Context() bot.Context {
	return c.c
}

func (c *Context) Error(e error) {
	if e == nil {
		return
	}

	c.done <- e
}

func (c *Context) Send(what interface{}, opts ...interface{}) {
	c.Error(c.c.Send(what, opts...))
}

func (c *Context) SendTo(to bot.Recipient, what interface{}, opts ...interface{}) {
	_, err := c.c.Bot().Send(to, what, opts...)
	c.Error(err)
}

func (c *Context) SendAlbum(a bot.Album, opts ...interface{}) {
	c.Error(c.c.SendAlbum(a, opts...))
}

func (c *Context) Reply(what interface{}, opts ...interface{}) {
	c.Error(c.c.Reply(what, opts...))
}

func (c *Context) Forward(msg bot.Editable, opts ...interface{}) {
	c.Error(c.c.Forward(msg, opts...))
}

func (c *Context) ForwardTo(to bot.Recipient, opts ...interface{}) {
	c.Error(c.c.ForwardTo(to, opts...))
}

func (c *Context) Edit(what interface{}, opts ...interface{}) {
	c.Error(c.c.Edit(what, opts...))
}

func (c *Context) EditCaption(caption string, opts ...interface{}) {
	c.Error(c.c.EditCaption(caption, opts...))
}

func (c *Context) EditOrSend(what interface{}, opts ...interface{}) {
	c.Error(c.c.EditOrSend(what, opts...))
}

func (c *Context) EditOrReply(what interface{}, opts ...interface{}) {
	c.Error(c.c.EditOrReply(what, opts...))
}

func (c *Context) Delete() {
	c.Error(c.c.Delete())
}

func (c *Context) Notify(action bot.ChatAction) {
	c.Error(c.c.Notify(action))
}

func (c *Context) Answer(resp *bot.QueryResponse) {
	c.Error(c.c.Answer(resp))
}

func (c *Context) Respond(resp ...*bot.CallbackResponse) {
	c.Error(c.c.Respond(resp...))
}
