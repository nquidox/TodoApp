package todoList

import (
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
}

func (t *TodoList) Create(title string) error {
	err := validateListOnCreate(title)
	if err != nil {
		return err
	}

	t.Uuid = uuid.New()
	t.Title = title
	t.AddedDate = time.Now()
	t.Order = 0

	err = DB.Create(t).Error
	if err != nil {
		return err
	}

	return nil
}

func (t *TodoList) GetAllLists() ([]TodoList, error) {
	var allLists []TodoList
	err := DB.Find(&allLists).Error
	if err != nil {
		return []TodoList{}, err
	}
	return allLists, nil
}

func (t *TodoList) Update(id uuid.UUID, title string) error {
	err := DB.Where("uuid = ?", id).First(t).Error
	if err != nil {
		return err
	}

	t.Title = title
	err = DB.Save(t).Error
	if err != nil {
		return err
	}
	return nil
}

func (t *TodoList) Delete(uuid uuid.UUID) error {
	err := DB.Where("uuid = ?", uuid).Delete(t).Error
	if err != nil {
		return err
	}
	return nil
}