package todoList

import (
	"errors"
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
	var err error

	t.Description = ""
	t.Completed = "false"
	t.Status = 0
	t.Priority = 1
	t.StartDate = time.Time{}
	t.Deadline = time.Time{}
	t.TaskId = uuid.New()
	t.Order = 0
	t.AddedDate = time.Now()

	err = DB.Create(t).Error
	if err != nil {
		return err
	}

	return nil
}

func (t *Task) Read(count, page int) ([]Task, error) {
	var tasks []Task

	result := DB.
		Where("todo_list_id = ?", t.TodoListId).
		Offset((page - 1) * count).
		Limit(count).
		Find(&tasks)

	if result.Error != nil {
		return nil, result.Error
	}

	if result.RowsAffected == 0 {
		return nil, errors.New("not found")
	}

	return tasks, nil
}

func (t *Task) Update() error {
	result := DB.
		Where("todo_list_id = ?", t.TodoListId).
		Where("task_id = ?", t.TaskId).
		Updates(t)

	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected == 0 {
		return errors.New("not found")
	}

	return nil
}

func (t *Task) Delete() error {
	result := DB.
		Where("todo_list_id = ?", t.TodoListId).
		Where("task_id = ?", t.TaskId).
		Delete(t)

	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected == 0 {
		return errors.New("not found")
	}
	return nil
}
