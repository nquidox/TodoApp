package db

import (
	"fmt"
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

func (db *DB) ReadManyRecords(model any) error {
	err := db.Connection.Find(model).Error
	if err != nil {
		return err
	}
	return nil
}

func (db *DB) ReadWithPagination(model any, params map[string]any) error {
	err := db.Connection.
		Where(fmt.Sprintf("%s = ?", params["field"])).
		Offset((params["page"]).(int) - 1*params["count"].(int)).
		Limit(params["count"].(int)).
		Find(model).
		Error

	if err != nil {
		return err
	}

	return nil
}

func (db *DB) UpdateRecord(model any, params map[string]any) error {
	err := db.Connection.
		Where(fmt.Sprintf("%s = ?", params["field"]), params[params["field"].(string)]).
		Updates(model).
		Error
	if err != nil {
		return err
	}
	return nil
}

func (db *DB) DeleteRecord(model any, params map[string]any) error {
	err := db.Connection.
		Where(fmt.Sprintf("%s = ?", params["field"]), params[params["field"].(string)]).
		Delete(model).
		Error
	if err != nil {
		return err
	}
	return nil
}
