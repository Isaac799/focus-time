// Package sqlite keeps persistence
// syntax: https://github.com/mattn/go-sqlite3/blob/master/_example/simple/simple.go
// driver: https://pkg.go.dev/modernc.org/sqlite
package sqlite

import "errors"

var (
	// ErrWindowNameRequired is a preflight check before window db operations
	ErrWindowNameRequired = errors.New("cannot do operations on a nameless window")
	// DBName is the filename of the sqlite file for standard operation
	DBName = ".focustime.db"
	// DBNameTest is the filename of the sqlite file for tests, distinct to allow for dropping
	DBNameTest = ".focustime_test.db"
)
