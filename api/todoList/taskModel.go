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
	AddedDate   time.Time `json:"addedDate"`
}

func (t *Task) Create() error {
	t.Description = ""
	t.Completed = "false"
	t.Status = 0
	t.Priority = 1
	t.StartDate = time.Time{}
	t.Deadline = time.Time{}
	t.TaskId = uuid.New()
	t.Order = 0
	t.AddedDate = time.Now()

	err := Worker.CreateRecord(t)
	if err != nil {
		return err
	}

	return nil
}

func (t *Task) Read(count, page int) ([]Task, error) {
	var tasks []Task
	params := map[string]any{"field": "todo_list_id", "count": count, "page": page}
	err := Worker.ReadWithPagination(&tasks, params)
	if err != nil {
		return nil, err
	}
	return tasks, nil
}

func (t *Task) Update() error {
	//TODO: implement a structure to pass several parameters to worker
	//"todo_list_id"
	//"task_id"
	return nil
}

func (t *Task) Delete() error {
	//TODO: implement a structure to pass several parameters to worker
	//"todo_list_id"
	//"task_id"
	return nil
}
