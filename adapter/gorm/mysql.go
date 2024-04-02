package gorm

import (
	"fmt"

	"goapp/config"
	"goapp/util/logger/gorm_logger"

	gosql "github.com/go-sql-driver/mysql"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func openMysql() (*gorm.DB, error) {

	cfg := &gosql.Config{
		Net:                  "tcp",
		Addr:                 fmt.Sprintf("%v:%v", config.Conf.Database.Mysql.Host, config.Conf.Database.Mysql.Port),
		DBName:               config.Conf.Database.Mysql.Database,
		User:                 config.Conf.Database.Mysql.Username,
		Passwd:               config.Conf.Database.Mysql.Password,
		AllowNativePasswords: true,
		ParseTime:            true,
	}

	return gorm.Open(mysql.Open(cfg.FormatDSN()), &gorm.Config{
		Logger: gorm_logger.New(),
	})

}
