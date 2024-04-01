package repository

import (
	"goapp/models"

	"github.com/segmentio/ksuid"
	"gorm.io/gorm"
)

func CreateSession(db *gorm.DB, session *models.Session) (*models.Session, error) {
	session.ID = ksuid.New().String()
	if err := db.Create(session).Error; err != nil {
		return nil, err
	}
	return session, nil
}
