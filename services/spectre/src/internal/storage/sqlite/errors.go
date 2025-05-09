package sqlite

import "fmt"

var (
	ErrCannotConnectSQLite = func(loc string, err error) error {
		return fmt.Errorf("cannot connect to SQLite db at %s: %v", loc, err)
	}
)
