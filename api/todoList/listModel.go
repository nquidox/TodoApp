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
	ListUuid  uuid.UUID `json:"-"`
	OwnerUuid uuid.UUID `json:"-"`
	Title     string    `json:"title" binding:"required"  extensions:"x-order=1"`
	Order     int       `json:"order" extensions:"x-order=2"`
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

func (c *createTodoList) Create(dbw dbWorker) error {
	list := TodoList{
		ListUuid:  c.ListUuid,
		Title:     c.Title,
		Order:     c.Order,
		OwnerUuid: c.OwnerUuid,
	}

	err := dbw.CreateRecord(&list)
	if err != nil {
		return err
	}

	return nil
}

func (r *readTodoList) GetAllLists(dbw dbWorker, aw authUser) ([]readTodoList, error) {
	var allLists []readTodoList
	params := map[string]any{"owner_uuid": aw.UserUUID}
	err := dbw.ReadManyRecords(TodoList{}, &allLists, params)
	if err != nil {
		return nil, err
	}
	return allLists, nil
}

func (c *createTodoList) Update(dbw dbWorker) error {
	params := map[string]any{"list_uuid": c.ListUuid}
	err := dbw.UpdateRecordSubmodel(TodoList{}, c, params)
	if err != nil {
		return err
	}
	return nil
}

func (t *TodoList) Delete(dbw dbWorker) error {
	params := map[string]any{"list_uuid": t.ListUuid}
	err := dbw.DeleteRecord(t, params)
	if err != nil {
		return err
	}
	return nil
}
