package sqlite

import "fmt"

var (
	errCannotConnectSQLite = func(loc string, err error) error {
		return fmt.Errorf("cannot connect to SQLite db at %s: %v", loc, err)
	}
	errLetterNotFound = func(id int) error {
		return fmt.Errorf("letter with id %d not found", id)
	}
	errCannotGetWithID = func(id int, err error) error {
		return fmt.Errorf("cannot get with id %d: %v", id, err)
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
	errWhenAddingLetter = func(body string, err error) error {
		return fmt.Errorf("cannot add letter with body %s: %v", body, err)
	}
	errWhenUpdateLetter = func(body string, err error) error {
		return fmt.Errorf("cannot update letter with body %s: %v", body, err)
	}
	errUserNotFoundByLogin = func(login string) error {
		return fmt.Errorf("user with login %s not found", login)
	}
	errUserNotFoundByID = func(id int) error {
		return fmt.Errorf("user with id %d not found", id)
	}
	errCannotGetUser = func(login string, err error) error {
		return fmt.Errorf("cannot get user with login %s: %v", login, err)
	}
	errCannotSaveUser = func(login string, err error) error {
		return fmt.Errorf("cannot save user with login %s: %v", login, err)
	}
	errCannotDeleteUser = func(id int, err error) error {
		return fmt.Errorf("cannot delete user with id %d: %v", id, err)
	}
	errCannotUpdateUser = func(id int, err error) error {
		return fmt.Errorf("cannot update user with id %d: %v", id, err)
	}
)
