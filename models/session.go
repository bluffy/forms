package models

import "github.com/bluffy/forms/adapter/gorm"

type Sessions []*Session
type Session struct {
	gorm.ModelUID
	BrowserAgent string
	UserID       string
	User         User
}

type SessionDtos []*SessionDto
type SessionDto struct {
	ID           string `json:"id"`
	BrowserAgent string `json:"browser_agent"`
}

func (o Session) ToDto() *SessionDto {
	return &SessionDto{
		ID:           o.ID,
		BrowserAgent: o.BrowserAgent,
	}
}
func (os Sessions) ToDto() SessionDtos {
	dtos := make([]*SessionDto, len(os))
	for i, o := range os {
		dtos[i] = o.ToDto()
	}
	return dtos
}
