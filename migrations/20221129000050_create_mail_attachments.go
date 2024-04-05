package migrations

import (
	"context"
	"database/sql"

	"github.com/pressly/goose/v3"
)

func init() {
	goose.AddMigrationContext(Up_20221129000050, Down_20221129000050)
}

func Up_20221129000050(ctx context.Context, txn *sql.Tx) error {

	sql := "missing dialect"
	switch dbType := GetType(); dbType {
	default:
		sql = `
			CREATE TABLE IF NOT EXISTS mail_attachments
			(
				id             CHAR(27)       NOT NULL,
				mail_id        CHAR(27) 	  NOT NULL,
				file_id        CHAR(27) 	  NOT NULL,
				created_at     TIMESTAMP      NOT NULL,
				updated_at     TIMESTAMP      NULL,
				deleted_at     TIMESTAMP      NULL,
				PRIMARY KEY (ID),
				FOREIGN KEY (file_id) REFERENCES files(id) ON DELETE CASCADE
				FOREIGN KEY (mail_id) REFERENCES mails(id) ON DELETE CASCADE
			);`
	}

	_, err := txn.ExecContext(ctx, sql)

	return err

}

func Down_20221129000050(ctx context.Context, txn *sql.Tx) error {
	_, err := txn.Exec("DROP TABLE IF EXISTS mail_attachments;")
	return err
}
