package todoList

import (
	"errors"
	"github.com/google/uuid"
	"gorm.io/gorm"
	"time"
)

type TodoList struct {
	gorm.Model `json:"-"`
	Uuid       uuid.UUID `json:"id"`
	Title      string    `json:"title"`
	AddedDate  time.Time `json:"addedDate"`
	Order      int       `json:"order"`
	OwnerUuid  uuid.UUID `json:"-"`
}

type Item struct {
	List TodoList `json:"item"`
}

func (t *TodoList) Create() error {
	t.Uuid = uuid.New()
	t.AddedDate = time.Now()
	t.Order = 0

	err := DB.Create(t).Error
	if err != nil {
		return err
	}

	return nil
}

func (t *TodoList) GetAllLists() ([]TodoList, error) {
	var allLists []TodoList
	err := DB.Find(&allLists).Error
	if err != nil {
		return nil, err
	}
	return allLists, nil
}

func (t *TodoList) Update() error {
	result := DB.Where("uuid = ?", t.Uuid).Updates(t)

	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected == 0 {
		return errors.New("list not found")
	}
	return nil
}

func (t *TodoList) Delete() error {
	result := DB.Where("uuid = ?", t.Uuid).Delete(t)

	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected == 0 {
		return errors.New("list not found")
	}
	return nil
}
