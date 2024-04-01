package db

import (
	"database/sql"
	"fmt"

	"goapp/config"

	"github.com/go-sql-driver/mysql"
)

func openMysql(conf *config.Config) (*sql.DB, error) {
	cfg := &mysql.Config{
		Net:                  "tcp",
		Addr:                 fmt.Sprintf("%v:%v", conf.Database.Host, conf.Database.Port),
		DBName:               conf.Database.Database,
		User:                 conf.Database.Username,
		Passwd:               conf.Database.Password,
		AllowNativePasswords: true,
		ParseTime:            true,
	}
	return sql.Open("mysql", cfg.FormatDSN())
}
