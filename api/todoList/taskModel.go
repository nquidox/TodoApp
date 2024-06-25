package todoList

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
	"time"
)

type Task struct {
	gorm.Model
	Description string    `json:"description"`
	Title       string    `json:"title"`
	Completed   bool      `json:"completed"`
	Status      int       `json:"status"`
	Priority    int       `json:"priority"`
	StartDate   time.Time `json:"startDate"`
	Deadline    time.Time `json:"deadline"`
	Id          uuid.UUID `json:"id"`
	TodoListId  uuid.UUID `json:"todoListId"`
	Order       int       `json:"order"`
	AddedDate   time.Time `json:"addedDate"`
}
