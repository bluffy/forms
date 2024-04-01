package gorm

import (
	"goapp/config"
	"goapp/util/logger/gorm_logger"

	log "github.com/sirupsen/logrus"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func openSqlite() (*gorm.DB, error) {

	log.Info("SQL File: " + config.Conf.Database.Path)

	return gorm.Open(sqlite.Open(config.Conf.Database.Path), &gorm.Config{
		Logger: gorm_logger.New(),
	})

}
