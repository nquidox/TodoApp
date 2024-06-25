package todoList

import "errors"

func validateListOnCreate(title string) error {
	if len(title) < 1 {
		return errors.New("list has no title")
	}
	if len(title) > 100 {
		return errors.New("title is too long")
	}
	return nil
}
