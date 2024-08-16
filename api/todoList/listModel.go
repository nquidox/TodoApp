package todoList

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
	"time"
	"todoApp/types"
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

func (c *createTodoList) Create(wrk types.DatabaseWorker) error {
	list := TodoList{
		ListUuid:  uuid.New(),
		Title:     c.Title,
		Order:     c.Order,
		OwnerUuid: uuid.Nil, //change to uuid from auth token
	}

	err := wrk.CreateRecord(&list)
	if err != nil {
		return err
	}

	return nil
}

func (r *readTodoList) GetAllLists(wrk types.DatabaseWorker) ([]readTodoList, error) {
	var allLists []readTodoList
	err := wrk.ReadManyRecords(TodoList{}, &allLists)
	if err != nil {
		return nil, err
	}
	return allLists, nil
}

func (c *createTodoList) Update(wrk types.DatabaseWorker) error {
	params := map[string]any{"list_uuid": c.ListUuid}
	err := wrk.UpdateRecordSubmodel(TodoList{}, c, params)
	if err != nil {
		return err
	}
	return nil
}

func (t *TodoList) Delete(wrk types.DatabaseWorker) error {
	params := map[string]any{"list_uuid": t.ListUuid}
	err := wrk.DeleteRecord(t, params)
	if err != nil {
		return err
	}
	return nil
}
