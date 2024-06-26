package todoList

import (
	"errors"
	"strconv"
)

func validateListOnCreate(title string) error {
	if len(title) < 1 {
		return errors.New("list has no title")
	}
	if len(title) > 100 {
		return errors.New("title is too long")
	}
	return nil
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
