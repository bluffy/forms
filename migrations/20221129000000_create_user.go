package migrations

import (
	"context"
	"database/sql"
	"goapp/util/tools"
	"os"

	"github.com/pressly/goose/v3"
	"github.com/segmentio/ksuid"
)

func init() {
	goose.AddMigrationContext(Up_20221129000000, Down_20221129000000)
}

func Up_20221129000000(ctx context.Context, txn *sql.Tx) error {

	id := ksuid.New().String()
	email := "dev@bluffy.de"
	password := "mgr"

	envAdminEmail := os.Getenv("INIT_ADMIN_EMAIL")
	envAdminPassword := os.Getenv("INIT_ADMIN_PASSWORD")

	if envAdminEmail != "" {
		email = envAdminEmail
	}
	if envAdminPassword != "" {
		password = envAdminPassword
	}

	password, err := tools.HashPassword(password)
	if err != nil {
		return err
	}

	sql := "missing dialect"
	switch dbType := GetType(); dbType {
	default:
		sql = `
			CREATE TABLE IF NOT EXISTS users
			(
				id             		 CHAR(27)     NOT NULL,
				email          		 VARCHAR(255) NOT NULL,
				password       		 VARCHAR(100) NOT NULL,
				is_admin       		 TINYINT(1)   NULL,
				first_name	   		 VARCHAR2(30) NULL,
				last_name	   		 VARCHAR2(30) NULL,
				newsletter	         TINYINT(1)   NULL,
				new_password_request TIMESTAMP 	  NULL,
				created_at     	     TIMESTAMP    NOT NULL,
				updated_at     	     TIMESTAMP    NULL,
				deleted_at     		 TIMESTAMP    NULL,
				PRIMARY KEY (ID),
				UNIQUE (email)
			);`
	}

	_, err = txn.ExecContext(ctx, sql)

	if err != nil {
		return err
	}

	sql = "missing dialect"
	switch dbType := GetType(); dbType {
	case "mysql":
		sql = `
		INSERT INTO users (id,email,password,is_admin,created_at, updated_at) 
		VALUES('` + id + `', '` + email + `','` + password + `', 1, NOW(), NOW());
		`
	case "sqlite":
		sql = `
		INSERT INTO users (id,email,password,is_admin,created_at, updated_at) 
		VALUES('` + id + `', '` + email + `','` + password + `', 1, TIME(),TIME());
		`
	}

	_, err = txn.ExecContext(ctx, sql)
	return err

}

func Down_20221129000000(ctx context.Context, txn *sql.Tx) error {
	_, err := txn.Exec("DROP TABLE IF EXISTS users;")
	return err
}
