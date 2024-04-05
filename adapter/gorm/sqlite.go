package gorm

import (
	"goapp/config"
	"goapp/util/logger/gorm_logger"

	log "github.com/sirupsen/logrus"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func openSqlite(conf *config.SqliteConf) (*gorm.DB, error) {

	log.Info("SQL File: " + conf.Path)

	return gorm.Open(sqlite.Open(conf.Path), &gorm.Config{
		Logger: gorm_logger.New(),
	})

}
