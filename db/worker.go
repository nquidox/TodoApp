package db

import (
	"fmt"
	log "github.com/sirupsen/logrus"
)

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

func (db *DB) ReadOneRecord(model any, params map[string]any) error {
	query := db.Connection

	for key, value := range params {
		query = query.Where(fmt.Sprintf("%s = ?", key), value)
	}

	err := query.First(model).Error
	if err != nil {
		return err
	}
	return nil
}

func (db *DB) ReadManyRecords(model any, submodel any) error {
	err := db.Connection.
		Model(model).
		Find(submodel).
		Error
	if err != nil {
		return err
	}
	return nil
}

func (db *DB) ReadWithPagination(model any, params map[string]any) error {
	query := db.Connection
	debugLogHeader := "Read with pagination"

	for key, value := range params {
		switch key {
		case "page":
			{
				query = query.Offset((params["page"]).(int) - 1*params["count"].(int))
				log.WithFields(log.Fields{
					"page": params["page"],
				}).Debug(debugLogHeader)
			}
		case "count":
			{
				query = query.Limit(params["count"].(int))
				log.WithFields(log.Fields{
					"count": params["count"],
				}).Debug(debugLogHeader)
			}
		default:
			query = query.Where(fmt.Sprintf("%s = ?", key), value)
			log.WithFields(log.Fields{
				key: value,
			}).Debug(debugLogHeader)
		}
	}
	err := query.Find(model).Error
	if err != nil {
		return err
	}
	return nil
}

func (db *DB) UpdateRecord(model any, params map[string]any) error {
	query := db.Connection

	for k, v := range params {
		query = query.Where(fmt.Sprintf("%s = ?", k), v)
	}
	query = query.Updates(model)

	err := query.Error
	if err != nil {
		return err
	}
	return nil
}

func (db *DB) UpdateRecordSubmodel(model any, submodel any, params map[string]any) error {
	query := db.Connection.Model(model)

	for k, v := range params {
		query = query.Where(fmt.Sprintf("%s = ?", k), v)
	}

	query = query.Updates(submodel)

	err := query.Error
	if err != nil {
		return err
	}
	return nil
}

func (db *DB) DeleteRecord(model any, params map[string]any) error {
	query := db.Connection

	for k, v := range params {
		query = query.Where(fmt.Sprintf("%s = ?", k), v)
	}
	err := query.Delete(model).Error
	if err != nil {
		return err
	}
	return nil
}
