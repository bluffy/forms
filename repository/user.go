package repository

import (
	"goapp/models"

	"github.com/segmentio/ksuid"
	"gorm.io/gorm"
)

func ReadUser(db *gorm.DB, id string) (*models.User, error) {
	user := &models.User{}

	if err := db.Where("id = ?", id).First(&user).Error; err != nil {
		return nil, err
	}
	return user, nil
}

func GetUserByEmail(db *gorm.DB, email string) (*models.User, error) {
	user := &models.User{}

	if err := db.Where("email = ?", email).First(&user).Error; err != nil {
		return nil, err
	}
	return user, nil
}

func ReadRegisterUser(db *gorm.DB, id string) (*models.RegisterUserForm, error) {
	user := &models.RegisterUserForm{}

	if err := db.Where("id = ?", id).First(&user).Error; err != nil {
		return nil, err
	}
	return user, nil
}

func CreateRegisterUser(db *gorm.DB, obj *models.RegisterUser) (*models.RegisterUser, error) {
	if obj.ID == "" {
		obj.ID = ksuid.New().String()
	}
	if err := db.Create(obj).Error; err != nil {
		return nil, err
	}
	return obj, nil
}
