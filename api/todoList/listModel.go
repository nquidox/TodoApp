package todoList

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
	"time"
)

type TodoList struct {
	gorm.Model `json:"-"`
	ListUuid   uuid.UUID  `json:"id"`
	Title      string     `json:"title"`
	Order      int        `json:"order"`
	OwnerUuid  uuid.UUID  `json:"-"`
	StartDate  *time.Time `json:"startDate"`
	EndDate    *time.Time `json:"endDate"`
	TextColor  string     `json:"textColor"`
	BgColor    string     `json:"backgroundColor"`
}

type createTodoList struct {
	ListUuid  uuid.UUID  `json:"-"`
	OwnerUuid uuid.UUID  `json:"-"`
	Title     string     `json:"title" binding:"required"  extensions:"x-order=1"`
	Order     int        `json:"order" extensions:"x-order=2"`
	StartDate *time.Time `json:"startDate"`
	EndDate   *time.Time `json:"endDate"`
	TextColor string     `json:"textColor"`
	BgColor   string     `json:"backgroundColor"`
}

type readTodoList struct {
	ListUuid  uuid.UUID  `json:"id"`
	Title     string     `json:"title"`
	AddedDate time.Time  `json:"addedDate" gorm:"column:created_at"`
	Order     int        `json:"order"`
	OwnerUuid uuid.UUID  `json:"-"`
	StartDate *time.Time `json:"startDate"`
	EndDate   *time.Time `json:"endDate"`
	TextColor string     `json:"textColor"`
	BgColor   string     `json:"backgroundColor"`
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
		StartDate: c.StartDate,
		EndDate:   c.EndDate,
		TextColor: c.TextColor,
		BgColor:   c.BgColor,
	}

	err := dbw.CreateRecord(&list)
	if err != nil {
		return err
	}

	return nil
}

func (r *readTodoList) GetAllLists(dbw dbWorker, aw authUser, order string) ([]readTodoList, error) {
	var allLists []readTodoList
	params := map[string]any{"owner_uuid": aw.UserUUID, "order": order, "sort_by": "created_at"}
	err := dbw.ReadManyRecords(TodoList{}, &allLists, params)
	if err != nil {
		return nil, err
	}
	return allLists, nil
}

func (c *createTodoList) Update(dbw dbWorker) error {
	params := map[string]any{"list_uuid": c.ListUuid, "owner_uuid": c.OwnerUuid}
	err := dbw.UpdateRecordSubmodel(TodoList{}, c, params)
	if err != nil {
		return err
	}
	return nil
}

func (t *TodoList) Delete(dbw dbWorker) error {
	params := map[string]any{"list_uuid": t.ListUuid, "owner_uuid": t.OwnerUuid}
	err := dbw.DeleteRecord(t, params)
	if err != nil {
		return err
	}
	return nil
}
