package todoList

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
	"time"
)

type TodoList struct {
	gorm.Model `json:"-"`
	ListUuid   uuid.UUID `json:"id"`
	Title      string    `json:"title"`
	Order      int       `json:"order"`
	OwnerUuid  uuid.UUID `json:"-"`
}

type createTodoList struct {
	ListUuid uuid.UUID `json:"-"`
	Title    string    `json:"title" binding:"required"  extensions:"x-order=1"`
	Order    int       `json:"order" extensions:"x-order=2"`
}

type readTodoList struct {
	ListUuid  uuid.UUID `json:"id"`
	Title     string    `json:"title"`
	AddedDate time.Time `json:"addedDate" gorm:"column:created_at"`
	Order     int       `json:"order"`
	OwnerUuid uuid.UUID `json:"-"`
}

type Item struct {
	List createTodoList `json:"item"`
}

func (c *createTodoList) Create() error {
	list := TodoList{
		ListUuid:  uuid.New(),
		Title:     c.Title,
		Order:     c.Order,
		OwnerUuid: uuid.Nil, //change to uuid from auth token
	}

	err := Worker.CreateRecord(&list)
	if err != nil {
		return err
	}

	return nil
}

func (r *readTodoList) GetAllLists() ([]readTodoList, error) {
	var allLists []readTodoList
	err := Worker.ReadManyRecords(TodoList{}, &allLists)
	if err != nil {
		return nil, err
	}
	return allLists, nil
}

func (c *createTodoList) Update() error {
	params := map[string]any{"list_uuid": c.ListUuid}
	err := Worker.UpdateRecordSubmodel(TodoList{}, c, params)
	if err != nil {
		return err
	}
	return nil
}

func (t *TodoList) Delete() error {
	params := map[string]any{"list_uuid": t.ListUuid}
	err := Worker.DeleteRecord(t, params)
	if err != nil {
		return err
	}
	return nil
}
