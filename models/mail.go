package models

import (
	"goapp/adapter/gorm"
	"goapp/app/service"
	"strings"
	"time"
)

type Mails []*Mail
type Mail struct {
	gorm.ModelUID
	Status       int
	UserID       string
	User         User
	Sender       string
	Recipient    string
	Cc           *string
	Bc           *string
	Subject      string
	Text         *string
	Html         *string
	SendAt       *time.Time
	Error        *string
	ErrorMessage *string
}

func (m *Mail) ToServiceMail() *service.Mail {

	ret := &service.Mail{
		To:      m.Recipient,
		From:    m.Sender,
		Subject: m.Subject,
	}
	if m.Cc != nil {
		ret.Cc = strings.Split(*m.Cc, ",")
	}
	if m.Bc != nil {
		ret.Bc = strings.Split(*m.Bc, ",")
	}
	if m.Text != nil {
		ret.Text = *m.Text
	}
	if m.Html != nil {
		ret.HTML = *m.Html
	}

	return ret

}
