package gorm

import (
	"errors"
	"goapp/config"
	"time"

	"gorm.io/gorm"
)

// gorm.Model definition
type ModelUID struct {
	ID        string `gorm:"primaryKey"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`
}

func New() (*gorm.DB, error) {
	if config.Conf.Database.Type == "mysql" {
		return openMysql()
	}
	return nil, errors.New("no database Connector found! wrong type? (mysql,sqlite)")

}
