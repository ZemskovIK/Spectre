package storage

import "fmt"

var (
	ErrLetterNotFound = func(id int) error {
		return fmt.Errorf("letter with id %d not found", id)
	}
)
