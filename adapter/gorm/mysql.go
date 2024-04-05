package gorm

import (
	"fmt"

	"goapp/config"
	"goapp/util/logger/gorm_logger"

	gosql "github.com/go-sql-driver/mysql"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func openMysql(conf *config.MysqlConf) (*gorm.DB, error) {

	cfg := &gosql.Config{
		Net:                  "tcp",
		Addr:                 fmt.Sprintf("%v:%v", conf.Host, conf.Port),
		DBName:               conf.Database,
		User:                 conf.Username,
		Passwd:               conf.Password,
		AllowNativePasswords: true,
		ParseTime:            true,
	}

	return gorm.Open(mysql.Open(cfg.FormatDSN()), &gorm.Config{
		Logger: gorm_logger.New(),
	})

}
