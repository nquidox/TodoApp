package db

import (
	"errors"
	"fmt"
	log "github.com/sirupsen/logrus"
	"gorm.io/gorm/clause"
	"time"
)

func (db *DB) InitTable(model any) error {
	err := db.Connection.AutoMigrate(model)
	if err != nil {
		return err
	}
	return nil
}

func (db *DB) CreateRecord(model any) error {
	result := db.Connection.Create(model)

	if result.Error != nil {
		return result.Error
	}
	return nil
}

func (db *DB) ReadOneRecord(model any, params map[string]any) error {
	query := db.Connection

	for key, value := range params {
		switch key {
		case "model":
			query = query.Model(params["model"])
		default:
			query = query.Where(fmt.Sprintf("%s = ?", key), value)
		}
	}

	result := query.First(model)

	if result.RowsAffected == 0 {
		return errors.New("404")
	}

	if result.Error != nil {
		return result.Error
	}
	return nil
}

func (db *DB) ReadRecordSubmodel(model any, submodel any, params map[string]any) error {
	query := db.Connection.Model(model)

	for key, value := range params {
		query = query.Where(fmt.Sprintf("%s = ?", key), value)
	}

	result := query.First(submodel)

	if result.RowsAffected == 0 {
		return errors.New("404")
	}

	if result.Error != nil {
		return result.Error
	}
	return nil
}

func (db *DB) ReadManyRecords(model any, submodel any, params map[string]any) error {
	query := db.Connection.Model(model)

	for key, value := range params {
		switch key {
		case "order":
			query = query.Order(clause.OrderByColumn{
				Column: clause.Column{Name: params["sort_by"].(string)},
				Desc:   params["order"].(string) == "desc",
			})

		case "sort_by":
			continue

		default:
			query = query.Where(fmt.Sprintf("%s = ?", key), value)
		}
	}

	result := query.Find(submodel)

	if result.RowsAffected == 0 {
		return errors.New("404")
	}

	if result.Error != nil {
		return result.Error
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

		case "order":
			query = query.Order(clause.OrderByColumn{
				Column: clause.Column{Name: params["sort_by"].(string)},
				Desc:   params["order"].(string) == "desc",
			})

		case "sort_by":
			continue

		default:
			query = query.Where(fmt.Sprintf("%s = ?", key), value)
			log.WithFields(log.Fields{
				key: value,
			}).Debug(debugLogHeader)
		}
	}
	result := query.Find(model)

	if result.RowsAffected == 0 {
		return errors.New("404")
	}

	if result.Error != nil {
		return result.Error
	}
	return nil
}

func (db *DB) UpdateRecord(model any, params map[string]any) error {
	query := db.Connection

	for k, v := range params {
		query = query.Where(fmt.Sprintf("%s = ?", k), v)
	}
	result := query.Updates(model)

	if result.RowsAffected == 0 {
		return errors.New("404")
	}

	if result.Error != nil {
		return result.Error
	}
	return nil
}

func (db *DB) UpdateRecordSubmodel(model any, submodel any, params map[string]any) error {
	query := db.Connection.Model(model)

	for k, v := range params {
		query = query.Where(fmt.Sprintf("%s = ?", k), v)
	}

	result := query.Select("*").Updates(submodel)

	if result.RowsAffected == 0 {
		return errors.New("404")
	}

	if result.Error != nil {
		return result.Error
	}
	return nil
}

func (db *DB) DeleteRecord(model any, params map[string]any) error {
	query := db.Connection

	for k, v := range params {
		query = query.Where(fmt.Sprintf("%s = ?", k), v)
	}

	result := query.Delete(model)

	if result.RowsAffected == 0 {
		return errors.New("404")
	}

	if result.Error != nil {
		return result.Error
	}
	return nil
}

func (db *DB) DeleteManyExceptOne(model any, params map[string]any) error {
	query := `UPDATE sessions SET deleted_at = ? WHERE user_uuid = ? AND token != ? AND deleted_at IS NULL`
	result := db.Connection.Exec(query, time.Now(), params["user_uuid"], params["token"])

	if result.RowsAffected == 0 {
		return errors.New("404")
	}

	if result.Error != nil {
		return result.Error
	}
	return nil
}
