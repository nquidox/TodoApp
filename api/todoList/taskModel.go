package todoList

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
	"time"
)

type Task struct {
	gorm.Model  `json:"-"`
	Description string    `json:"description"`
	Title       string    `json:"title"`
	Completed   string    `json:"completed"`
	Status      int       `json:"status"`
	Priority    int       `json:"priority"`
	StartDate   time.Time `json:"startDate"`
	Deadline    time.Time `json:"deadline"`
	TaskId      uuid.UUID `json:"id"`
	TodoListId  uuid.UUID `json:"-"`
	Order       int       `json:"order"`
	AddedDate   time.Time `json:"addedDate" gorm:"column:created_at"`
}

type createTask struct {
	Title       string    `json:"title" extensions:"x-order=1"`
	Description string    `json:"description" extensions:"x-order=2"`
	Status      int       `json:"status" extensions:"x-order=3"`
	Priority    int       `json:"priority" extensions:"x-order=4"`
	Order       int       `json:"order" extensions:"x-order=5"`
	StartDate   string    `json:"startDate" extensions:"x-order=6"`
	Deadline    string    `json:"deadline" extensions:"x-order=7"`
	TodoListId  uuid.UUID `json:"-"`
	TaskId      uuid.UUID `json:"-"`
}

func (t *Task) Create() error {
	err := Worker.CreateRecord(t)
	if err != nil {
		return err
	}

	return nil
}

func (t *Task) Read(count, page int) ([]Task, error) {
	var tasks []Task
	params := map[string]any{"todo_list_id": t.TodoListId, "count": count, "page": page}
	err := Worker.ReadWithPagination(&tasks, params)
	if err != nil {
		return nil, err
	}
	return tasks, nil
}

func (c *createTask) Update() error {
	params := map[string]any{"todo_list_id": c.TodoListId, "task_id": c.TaskId}
	err := Worker.UpdateRecordSubmodel(Task{}, c, params)
	if err != nil {
		return err
	}
	return nil
}

func (t *Task) Delete() error {
	params := map[string]any{"todo_list_id": t.TodoListId, "task_id": t.TaskId}
	err := Worker.DeleteRecord(t, params)
	if err != nil {
		return err
	}
	return nil
}
