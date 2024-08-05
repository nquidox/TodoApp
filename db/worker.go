package db

import (
	"fmt"
	log "github.com/sirupsen/logrus"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"todoApp/config"
)

func ConnectToDB(c *config.Config) *gorm.DB {
	var level logger.LogLevel
	var err error

	switch c.Config.DBLogLevel {
	case "silent":
		level = logger.Silent
	case "error":
		level = logger.Error
	case "warn":
		level = logger.Warn
	default:
		{
			level = logger.Info
			log.Warning("Error parsing database level, using default: Info")
		}
	}

	dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		c.Config.Host, c.Config.Port, c.Config.User, c.Config.Password, c.Config.Dbname, c.Config.Sslmode)

	DB, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger:                 logger.Default.LogMode(level),
		SkipDefaultTransaction: true,
	})

	if err != nil {
		log.Fatal(err)
	}
	return DB
}
