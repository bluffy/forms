package models

import (
	"goapp/adapter/gorm"
	"time"
)

type Mails []*Mail
type Mail struct {
	gorm.ModelUID
	Status    int
	UserID    string
	User      User
	Sender    string
	Recipient string
	ReplyTo   *string
	Cc        *string
	Bc        *string
	Subject   string
	Text      *string
	Html      *string
	SendAt    *time.Time
	Error     *string
}
