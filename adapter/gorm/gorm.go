package gorm

import (
	"fmt"
	"time"

	"github.com/bluffy/forms/config"

	"github.com/bluffy/forms/util/logger/gorm_logger"

	gosql "github.com/go-sql-driver/mysql"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

// gorm.Model definition
type ModelUID struct {
	ID        string `gorm:"primaryKey"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`
}

func New(conf *config.Config) (*gorm.DB, error) {

	cfg := &gosql.Config{
		Net:                  "tcp",
		Addr:                 fmt.Sprintf("%v:%v", conf.Database.Host, conf.Database.Port),
		DBName:               conf.Database.Database,
		User:                 conf.Database.Username,
		Passwd:               conf.Database.Password,
		AllowNativePasswords: true,
		ParseTime:            true,
	}

	return gorm.Open(mysql.Open(cfg.FormatDSN()), &gorm.Config{
		Logger: gorm_logger.New(),
	})

}
