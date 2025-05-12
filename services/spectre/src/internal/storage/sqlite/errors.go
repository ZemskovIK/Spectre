package sqlite

import "fmt"

var (
	ErrCannotConnectSQLite = func(loc string, err error) error {
		return fmt.Errorf("cannot connect to SQLite db at %s: %v", loc, err)
	}
	ErrLetterNotFound = func(id int) error {
		return fmt.Errorf("letter with id %d not found", id)
	}
	ErrCannotGetLetter = func(id int, err error) error {
		return fmt.Errorf("cannot get letter with id %d: %v", id, err)
	}
	ErrWnenFetchAuthor = func(name string, err error) error {
		return fmt.Errorf("cannot fetch author %s: %v", name, err)
	}
	ErrWhenAddingAuthor = func(name string, err error) error {
		return fmt.Errorf("cannot add author %s: %v", name, err)
	}
	ErrWhenAddingLetter = func(title string, err error) error {
		return fmt.Errorf("cannot add letter with title %s: %v", title, err)
	}
)
