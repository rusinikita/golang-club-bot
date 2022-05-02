package request

import (
	"strconv"
	"strings"

	bot "gopkg.in/telebot.v3"
)

type Request struct {
	RecordID       string
	UserID         string
	Name           string
	Username       string
	RepositoryLink string `json:",omitempty"`
	Status         Status
	DeclineMessage string
	Send           bool // is notification send
}

func (r Request) Recipient() string {
	return r.UserID
}

type Status string

const (
	None     Status = ""
	Accepted Status = "Accepted"
	Declined Status = "Declined"
)

func NewID(u *bot.User) Request {
	return Request{UserID: strconv.FormatInt(u.ID, 10)}
}

func New(u *bot.User) Request {
	return Request{
		UserID:   strconv.FormatInt(u.ID, 10),
		Name:     strings.Join([]string{u.FirstName, u.LastName}, " "),
		Username: u.Username,
	}
}
