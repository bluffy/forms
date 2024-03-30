package repository

import (
	"github.com/bluffy/forms/models"
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
