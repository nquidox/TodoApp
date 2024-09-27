package todoList

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
	"time"
)

type Task struct {
	gorm.Model   `json:"-"`
	Description  string     `json:"description"`
	Title        string     `json:"title"`
	Completed    string     `json:"completed"`
	Status       int        `json:"status"`
	Priority     int        `json:"priority"`
	StartDate    *time.Time `json:"startDate"`
	Deadline     *time.Time `json:"deadline"`
	TaskUUID     uuid.UUID  `json:"id" gorm:"index"`
	TodoListUUID uuid.UUID  `json:"-" gorm:"index"`
	Order        int        `json:"order"`
	AddedDate    time.Time  `json:"addedDate" gorm:"column:created_at; autoCreateTime"`
	OwnerUUID    uuid.UUID  `json:"-" gorm:"index"`
}

type createTask struct {
	Title        string     `json:"title" extensions:"x-order=1"`
	Description  string     `json:"description" extensions:"x-order=2"`
	Status       int        `json:"status" extensions:"x-order=3"`
	Priority     int        `json:"priority" extensions:"x-order=4"`
	Order        int        `json:"order" extensions:"x-order=5"`
	StartDate    *time.Time `json:"startDate" extensions:"x-order=6"`
	Deadline     *time.Time `json:"deadline" extensions:"x-order=7"`
	TodoListUUID uuid.UUID  `json:"-"`
	TaskUUID     uuid.UUID  `json:"-"`
	OwnerUUID    uuid.UUID  `json:"-"`
}

func (t *Task) Create(dbw dbWorker) error {
	err := dbw.CreateRecord(t)
	if err != nil {
		return err
	}
	return nil
}

func (t *Task) Read(dbw dbWorker, order string, count, page int) ([]Task, error) {
	var tasks []Task
	params := map[string]any{
		"todo_list_uuid": t.TodoListUUID,
		"count":          count,
		"page":           page,
		"owner_uuid":     t.OwnerUUID,
		"order":          order,
		"sort_by":        "created_at",
	}
	err := dbw.ReadWithPagination(&tasks, params)
	if err != nil {
		return nil, err
	}
	return tasks, nil
}

func (c *createTask) Update(dbw dbWorker) error {
	params := map[string]any{
		"todo_list_uuid": c.TodoListUUID,
		"task_uuid":      c.TaskUUID,
		"owner_uuid":     c.OwnerUUID,
	}
	err := dbw.UpdateRecordSubmodel(Task{}, c, params)
	if err != nil {
		return err
	}
	return nil
}

func (t *Task) Delete(dbw dbWorker) error {
	params := map[string]any{
		"todo_list_uuid": t.TodoListUUID,
		"task_uuid":      t.TaskUUID,
		"owner_uuid":     t.OwnerUUID,
	}
	err := dbw.DeleteRecord(t, params)
	if err != nil {
		return err
	}
	return nil
}
