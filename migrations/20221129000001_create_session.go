package migrations

import (
	"context"
	"database/sql"
	"goapp/config"

	"github.com/pressly/goose/v3"
)

func init() {
	goose.AddMigrationContext(Up_20221129000001, Down_20221129000001)
}

func Up_20221129000001(ctx context.Context, tx *sql.Tx) error {

	sql := "missing dialect"
	switch dbType := config.Conf.Database.Type; dbType {

	case "mysql":
		sql = `
        CREATE TABLE session (
            key       CHAR(64) NOT NULL,
            data      BYTEA,
            expiry    INTEGER NOT NULL,
            PRIMARY KEY (key)
        );
        `
	case "sqlite":
		sql = `
        CREATE TABLE session (
            key       CHAR(64) NOT NULL,
            data      BYTEA,
            expiry    INTEGER NOT NULL,
            PRIMARY KEY (key)
        );
        `
		/*

				case "mysql":
					sql = `
			        CREATE TABLE IF NOT EXISTS sessions
			        (
			            id             CHAR(27)     NOT NULL,
			            user_id        CHAR(27)     NOT NULL,
			            browser_agent  VARCHAR(1000) NULL,
			            created_at     TIMESTAMP    NOT NULL,
			            updated_at     TIMESTAMP    NULL,
			            deleted_at     TIMESTAMP    NULL,
			            PRIMARY KEY (ID),
			            FOREIGN KEY (user_id)
			                REFERENCES users(id)
			                ON DELETE CASCADE
			        );
			        `
				case "sqlite":
					sql = `
			        CREATE TABLE IF NOT EXISTS sessions
			        (
			            id             CHAR(27)     NOT NULL,
			            user_id        CHAR(27)     NOT NULL,
			            browser_agent  VARCHAR(1000) NULL,
			            created_at     TIMESTAMP    NOT NULL,
			            updated_at     TIMESTAMP    NULL,
			            deleted_at     TIMESTAMP    NULL,
			            PRIMARY KEY (ID),
			            FOREIGN KEY (user_id)
			                REFERENCES users(id)
			                ON DELETE CASCADE
			        );
			        `
		*/
	}

	_, err := tx.ExecContext(ctx, sql)

	if err != nil {
		return err
	}

	return err

}

func Down_20221129000001(ctx context.Context, txn *sql.Tx) error {
	_, err := txn.ExecContext(ctx, "DROP TABLE IF EXISTS sessions;")
	return err
}
