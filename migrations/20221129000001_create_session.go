package migrations

import (
	"database/sql"
	"goapp/config"

	"github.com/pressly/goose/v3"
)

func init() {
	goose.AddMigration(Up_20221129000001, Down_20221129000001)
}

func Up_20221129000001(txn *sql.Tx) error {

	sql := "missing dialect"
	switch dbType := config.Conf.Database.Type; dbType {
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
	}

	_, err := txn.Exec(sql)

	if err != nil {
		return err
	}

	return err

}

func Down_20221129000001(txn *sql.Tx) error {
	_, err := txn.Exec("DROP TABLE IF EXISTS sessions;")
	return err
}
