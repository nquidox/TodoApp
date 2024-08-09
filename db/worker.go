package db

import (
	"fmt"
	log "github.com/sirupsen/logrus"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"todoApp/config"
)

type DatabaseWorker interface {
	InitTable(model any) error
	CreateRecord(model any) error
	ReadOneRecord(model any, field string, value any) error
	ReadManyRecords(model any) error
	UpdateRecord(model any, field string, value any) error
	DeleteRecord(model any, field string, value any) error
}

type DB struct {
	Connection *gorm.DB
}

func (db *DB) InitTable(model any) error {
	err := db.Connection.AutoMigrate(model)
	if err != nil {
		return err
	}
	return nil
}

func (db *DB) CreateRecord(model any) error {
	err := db.Connection.Create(model).Error
	if err != nil {
		return err
	}
	return nil
}

func (db *DB) ReadOneRecord(model any, field string, value any) error {
	err := db.Connection.
		Where(fmt.Sprintf("%s = ?", field), value).
		First(model).
		Error
	if err != nil {
		return err
	}
	return nil
}

func (db *DB) ReadManyRecords(model any) error {
	return nil
}

func (db *DB) UpdateRecord(model any, field string, value any) error {
	err := db.Connection.
		Where(fmt.Sprintf("%s = ?", field), value).
		Updates(model).
		Error
	if err != nil {
		return err
	}
	return nil
}

func (db *DB) DeleteRecord(model any, field string, value any) error {
	err := db.Connection.
		Where(fmt.Sprintf("%s = ?", field), value).
		Delete(model).
		Error
	if err != nil {
		return err
	}
	return nil
}

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
