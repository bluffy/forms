package gorm

import (
	"goapp/config"
	"goapp/util/logger/gorm_logger"
	"log"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func openSqlite() (*gorm.DB, error) {

	log.Print(config.Conf.Database.Path)

	return gorm.Open(sqlite.Open(config.Conf.Database.Path), &gorm.Config{
		Logger: gorm_logger.New(),
	})

}
