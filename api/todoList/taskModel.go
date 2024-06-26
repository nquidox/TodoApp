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
	Completed   bool      `json:"completed"`
	Status      int       `json:"status"`
	Priority    int       `json:"priority"`
	StartDate   time.Time `json:"startDate"`
	Deadline    time.Time `json:"deadline"`
	TaskId      uuid.UUID `json:"id"`
	TodoListId  uuid.UUID `json:"-"`
	Order       int       `json:"order"`
	AddedDate   time.Time `json:"addedDate"`
}

func (t *Task) Create(listId uuid.UUID, title string) error {
	var err error

	t.Description = ""
	t.Title = title
	t.Completed = false
	t.Status = 0
	t.Priority = 1
	t.StartDate = time.Time{}
	t.Deadline = time.Time{}
	t.TaskId = uuid.New()
	t.TodoListId = listId
	t.Order = 0
	t.AddedDate = time.Now()

	err = DB.Create(t).Error
	if err != nil {
		return err
	}

	return nil
}

func (t *Task) Read(listId uuid.UUID, count, page int) ([]Task, error) {
	var err error
	var tasks []Task

	err = DB.
		Where("todo_list_id = ?", listId).
		Offset((page - 1) * count).
		Limit(count).
		Find(&tasks).
		Error

	if err != nil {
		return []Task{}, err
	}
	return tasks, nil
}

func (t *Task) Update(listId, taskId uuid.UUID) (Task, error) {
	var err error
	var list TodoList
	var currentTask Task

	err = DB.Where("uuid = ?", listId).First(&list).Error
	if err != nil {
		return Task{}, err
	}

	err = DB.Where("task_id = ?", taskId).First(&currentTask).Error
	if err != nil {
		return Task{}, err
	}

	err = DB.Model(&currentTask).Updates(t).Error
	if err != nil {
		return Task{}, err
	}

	return currentTask, nil
}

func (t *Task) Delete(listId uuid.UUID, taskId uuid.UUID) error {
	var err error
	var list TodoList

	err = DB.Where("uuid = ?", listId).First(&list).Error
	if err != nil {
		return err
	}

	err = DB.Where("task_id = ?", taskId).Delete(t).Error
	if err != nil {
		return err
	}
	return nil
}
