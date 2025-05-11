package db

import "errors"

var (
	// ErrWindowNameRequired is a preflight check before window db operations
	ErrWindowNameRequired = errors.New("cannot do operations on a nameless window")
	// ErrZeroID is used when an ID is expected, e.g. a read did not scan an ID
	ErrZeroID = errors.New("failed to read record id was zero")
	// ErrMissingValue is used when a value is missing, e.g. a read did not scan
	ErrMissingValue = errors.New("missing an expected value")
	// ErrUnexpectedValue is used when a value is not met, e.g. an update was incomplete
	ErrUnexpectedValue = errors.New("expected value was not met")
	// DBName is the filename of the sqlite file for standard operation
	DBName = ".focustime.db"
	// DBNameTest is the filename of the sqlite file for tests, distinct to allow for dropping
	DBNameTest = ".focustime_test.db"
)
