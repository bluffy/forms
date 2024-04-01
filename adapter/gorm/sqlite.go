package gorm

import (
	"goapp/util/logger/gorm_logger"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func openSqlite() (*gorm.DB, error) {

	return gorm.Open(sqlite.Open("database.db"), &gorm.Config{
		Logger: gorm_logger.New(),
	})

}
