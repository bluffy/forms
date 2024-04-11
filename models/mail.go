package models

import (
	"goapp/adapter/gorm"

	"time"
)

const (
	SEND_STATUS_NEW     int = 0
	SEND_STATUS_WAITING int = 1
	SEND_STATUS_SENT    int = 2
	SEND_STATUS_ERROR   int = 9
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
