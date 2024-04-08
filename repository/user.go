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

func ReadRegisterUser(db *gorm.DB, id string) (*models.RegisterUser, error) {
	user := &models.RegisterUser{}

	if err := db.Where("id = ?", id).First(&user).Error; err != nil {
		return nil, err
	}
	return user, nil
}

func DeleteRegisterUser(db *gorm.DB, id string) error {
	user := &models.RegisterUser{}
	if err := db.Where("id = ?", id).Delete(&user).Error; err != nil {
		return err
	}
	return nil
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
func CreateUser(db *gorm.DB, obj *models.User) (*models.User, error) {
	if obj.ID == "" {
		obj.ID = ksuid.New().String()
	}
	if err := db.Create(obj).Error; err != nil {
		return nil, err
	}
	return obj, nil
}

func CountAllUsers(db *gorm.DB) int64 {
	var count int64
	db.Model(&models.User{}).Where("name = ?", "jinzhu").Count(&count)

	return count
}
