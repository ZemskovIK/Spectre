package sqlite

import "fmt"

var (
	errCannotConnectSQLite = func(loc string, err error) error {
		return fmt.Errorf("cannot connect to SQLite db at %s: %v", loc, err)
	}
	errLetterNotFound = func(id int) error {
		return fmt.Errorf("letter with id %d not found", id)
	}
	errCannotGetLetter = func(id int, err error) error {
		return fmt.Errorf("cannot get letter with id %d: %v", id, err)
	}
	errCannotDeleteLetter = func(id int, err error) error {
		return fmt.Errorf("cannot delete letter with id %d: %w", id, err)
	}
	errCannotFetchRows = func(id int, err error) error {
		return fmt.Errorf("cannot fetch rows affected for id %d: %w", id, err)
	}
	errWnenFetchAuthor = func(name string, err error) error {
		return fmt.Errorf("cannot fetch author %s: %v", name, err)
	}
	errWhenAddingAuthor = func(name string, err error) error {
		return fmt.Errorf("cannot add author %s: %v", name, err)
	}
	errWhenAddingLetter = func(title string, err error) error {
		return fmt.Errorf("cannot add letter with title %s: %v", title, err)
	}
)
