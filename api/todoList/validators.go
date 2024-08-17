package todoList

import (
	"errors"
	"strconv"
	"time"
)

func (c *createTodoList) validateTitle() error {
	return validateTitle(c.Title, "list")
}

func (c *createTask) validateTitle() error {
	return validateTitle(c.Title, "task")
}

func (t *Task) validateTitle() error {
	return validateTitle(t.Title, "task")
}

func validateTime(timeString string) time.Time {
	t, err := time.Parse("02-01-2006 15:04:05", timeString)
	if err != nil {
		return time.Now()
	}
	return t
}

func validateQueryInt(queryValue string, defaultValue int) int {
	i, err := strconv.Atoi(queryValue)
	if err != nil {
		return defaultValue
	}
	if i < 0 {
		return defaultValue
	}
	return i
}

func validateTitle(title, fieldName string) error {
	if len(title) < 1 {
		return errors.New(fieldName + " has no title")
	}
	if len(title) > 100 {
		return errors.New("title is too long (MAX=100)")
	}
	return nil
}
