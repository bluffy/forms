package repository

import (
	"goapp/models"

	"github.com/segmentio/ksuid"
	"gorm.io/gorm"
)

func ReadMail(db *gorm.DB, id string) (*models.Mail, error) {
	obj := &models.Mail{}

	if err := db.Where("id = ?", id).First(&obj).Error; err != nil {
		return nil, err
	}
	return obj, nil
}
func CreateMail(db *gorm.DB, obj *models.Mail) (*models.Mail, error) {
	if obj.ID == "" {
		obj.ID = ksuid.New().String()
	}
	if err := db.Create(obj).Error; err != nil {
		return nil, err
	}
	return obj, nil
}
