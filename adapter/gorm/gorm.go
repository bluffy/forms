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

func New(conf *config.Database) (*gorm.DB, error) {
	if conf.Type == "mysql" {
		return openMysql(&conf.Mysql)
	}
	if conf.Type == "sqlite" {
		return openSqlite(&conf.Sqlite)
	}
	return nil, errors.New("no database Connector found! wrong type? (mysql,sqlite)")

}
