package todoList

import (
	"errors"
	"fmt"
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
	minChars := 1
	maxChars := 1000
	if len([]rune(title)) < minChars {
		return errors.New(fmt.Sprintf("%s has to be at least %d characters long.", fieldName, minChars))
	}
	if len([]rune(title)) > maxChars {
		return errors.New(fmt.Sprintf("%s is too long (MAX=%d)", fieldName, maxChars))
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
