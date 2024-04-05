package migrations

import (
	"log"
	"os"
)

func GetType() string {
	dbType := os.Getenv("BL_MIGRATE_DATABASE_TYPE")
	if dbType == "" {
		log.Fatal("No Database Type Found")
	}
	return os.Getenv("BL_MIGRATE_DATABASE_TYPE")
}
