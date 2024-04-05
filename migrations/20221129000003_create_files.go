package migrations

import (
	"context"
	"database/sql"

	"github.com/pressly/goose/v3"
)

func init() {
	goose.AddMigrationContext(Up_20221129000003, Down_20221129000003)
}

func Up_20221129000003(ctx context.Context, txn *sql.Tx) error {

	sql := "missing dialect"
	switch dbType := GetType(); dbType {
	default:
		sql = `
			CREATE TABLE IF NOT EXISTS files
			(
				id             CHAR(27)       NOT NULL,
				user_id        CHAR(27)       NULL,
				path           TEXT           NOT NULL,
				filename       VARCHAR(1000)  NOT NULL,
k				mime_type      VARCHAR(100)   NOT NULL,
				type           TINYINT        NOT NULL,
				size           INT            NOT NULL,
				height         INT            NULL,
				width          INT            NULL,
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

func Down_20221129000003(ctx context.Context, txn *sql.Tx) error {
	_, err := txn.Exec("DROP TABLE IF EXISTS files;")
	return err
}
