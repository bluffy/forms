package migrations

import (
	"context"
	"database/sql"

	"github.com/pressly/goose/v3"
)

func init() {
	goose.AddMigrationContext(Up_20221129000040, Down_20221129000040)
}

func Up_20221129000040(ctx context.Context, txn *sql.Tx) error {

	sql := "missing dialect"
	switch dbType := GetType(); dbType {
	default:
		sql = `
			CREATE TABLE IF NOT EXISTS mails
			(
				id             CHAR(27)       NOT NULL,
				user_id        CHAR(27)       NULL,
				status         tinyint(27)    NOT NULL,
				sender         VARCHAR(320)   NOT NULL,
				recipient      VARCHAR(320)   NOT NULL,
				cc             TEXT  		  NULL,
				bc             TEXT  		  NULL,
				subject        MEDIUMTEXT     NULL,
				text           MEDIUMTEXT     NULL,
				html           MEDIUMTEXT     NULL,
				send_at        TIMESTAMP      NULL,
				error_message   TEXT           NULL,
				error          TEXT           NULL,
				locale         VARCHAR(5)     NULL,
				created_at     TIMESTAMP      NOT NULL,
				updated_at     TIMESTAMP      NULL,
				deleted_at     TIMESTAMP      NULL,
				PRIMARY KEY (ID),
				FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE SET NULL
			);`
	}

	_, err := txn.ExecContext(ctx, sql)

	return err

}

func Down_20221129000040(ctx context.Context, txn *sql.Tx) error {
	_, err := txn.Exec("DROP TABLE IF EXISTS mails;")
	return err
}
