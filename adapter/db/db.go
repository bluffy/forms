package db

import (
	"database/sql"
	"errors"

	"goapp/config"
)

func New(conf *config.Config) (*sql.DB, error) {
	if conf.Database.Type == "mysql" {
		return openMysql(conf)
	}
	return nil, errors.New("no database Connector found! wrong type? (mysql,sqlite)")
	/*
		if conf.Db.Type == 2 {
			return open(conf)
		} else {
			return open(conf)
		}
	*/
}
