package todoList

import (
	"errors"
	"strconv"
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

func validateOrder(order string) string {
	switch order {
	case "asc":
		return "asc"
	default:
		return "desc"
	}
}
