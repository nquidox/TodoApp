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
	OwnerUuid  uuid.UUID `json:"-"`
}

type Item struct {
	List TodoList `json:"item"`
}

func (t *TodoList) Create() error {
	t.Uuid = uuid.New()
	t.AddedDate = time.Now()
	t.Order = 0

	err := Worker.CreateRecord(t)
	if err != nil {
		return err
	}

	return nil
}

func (t *TodoList) GetAllLists() ([]TodoList, error) {
	var allLists []TodoList
	err := Worker.ReadManyRecords(&allLists)
	if err != nil {
		return nil, err
	}
	return allLists, nil
}

func (t *TodoList) Update() error {
	params := map[string]any{"field": "uuid", "uuid": t.Uuid}
	err := Worker.UpdateRecord(t, params)
	if err != nil {
		return err
	}
	return nil
}

func (t *TodoList) Delete() error {
	params := map[string]any{"field": "uuid", "uuid": t.Uuid}
	err := Worker.DeleteRecord(t, params)
	if err != nil {
		return err
	}
	return nil
}
