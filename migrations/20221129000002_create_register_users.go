package migrations

import (
	"context"
	"database/sql"

	"github.com/pressly/goose/v3"
)

func init() {
	goose.AddMigrationContext(Up_20221129000002, Down_20221129000002)
}

func Up_20221129000002(ctx context.Context, txn *sql.Tx) error {

	sql := "missing dialect"
	switch dbType := GetType(); dbType {
	default:
		sql = `
			CREATE TABLE IF NOT EXISTS register_users
			(
				id             CHAR(27)     NOT NULL,
				email          VARCHAR(255) NOT NULL,
				password       VARCHAR(100) NOT NULL,
				first_name	   VARCHAR2(30) NULL,
				last_name	   VARCHAR2(30) NULL,
				newsletter	   TINYINT(1)   NULL,
				created_at     TIMESTAMP    NOT NULL,
				updated_at     TIMESTAMP    NULL,
				deleted_at     TIMESTAMP    NULL,
				PRIMARY KEY (ID)
			);`
	}

	_, err := txn.ExecContext(ctx, sql)

	return err

}

func Down_20221129000002(ctx context.Context, txn *sql.Tx) error {
	_, err := txn.Exec("DROP TABLE IF EXISTS register_users;")
	return err
}
