package migrations

import (
	"os"

	"github.com/sirupsen/logrus"
)

func GetType() string {
	dbType := os.Getenv("BL_MIGRATE_DATABASE_TYPE")
	if dbType == "" {
		logrus.Fatal("No Database Type Found")
	}
	return os.Getenv("BL_MIGRATE_DATABASE_TYPE")
}
