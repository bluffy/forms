package migrations

import (
	"database/sql"
	"goapp/config"
	"goapp/util/tools"

	"github.com/pressly/goose/v3"
	"github.com/segmentio/ksuid"
)

func init() {
	goose.AddMigration(Up_20221129000000, Down_20221129000000)
}

func Up_20221129000000(txn *sql.Tx) error {
	id := ksuid.New().String()

	email := "dev@bluffy.de"
	password, err := tools.HashPassword("mgr")

	if err != nil {
		return err
	}
	sql := ""
	switch dbType := config.Conf.Database.Type; dbType {
	case "mysql":
		sql = `
		CREATE TABLE IF NOT EXISTS users
		(
			id             CHAR(27)     NOT NULL,
			email          VARCHAR(255) NOT NULL,
			password       VARCHAR(255) NOT NULL,
			is_admin       TINYINT(1)   NULL,
			created_at     TIMESTAMP    NOT NULL,
			updated_at     TIMESTAMP    NULL,
			deleted_at     TIMESTAMP    NULL,
			PRIMARY KEY (ID)
		);`
	case "sqlite":
		sql = `
		CREATE TABLE IF NOT EXISTS users
		(
			id             CHAR(27)     NOT NULL,
			email          VARCHAR(255) NOT NULL,
			password       VARCHAR(255) NOT NULL,
			is_admin       TINYINT(1)   NULL,
			created_at     TIMESTAMP    NOT NULL,
			updated_at     TIMESTAMP    NULL,
			deleted_at     TIMESTAMP    NULL,
			PRIMARY KEY (ID)
		);`
	}

	_, err = txn.Exec(sql)

	if err != nil {
		return err
	}

	sql = ""
	switch dbType := config.Conf.Database.Type; dbType {
	case "mysql":
		sql = `
		INSERT INTO users (id,email,password,is_admin,created_at) 
		VALUES('` + id + `', '` + email + `','` + password + `', 1, NOW());
		`
	case "sqlite":
		sql = `
		INSERT INTO users (id,email,password,is_admin,created_at) 
		VALUES('` + id + `', '` + email + `','` + password + `', 1, TIME());
		`
	}

	_, err = txn.Exec(sql)
	return err

}

func Down_20221129000000(txn *sql.Tx) error {
	_, err := txn.Exec("DROP TABLE IF EXISTS users;")
	return err
}
