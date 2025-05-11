package db

import (
	"database/sql"
	"os"
	"path/filepath"
)

// Database stores the database connection and name
type Database struct {
	filename string
	DB       *sql.DB
	cache    cache
}

func connect(filename string) (*Database, error) {
	sqlite := Database{
		filename: "",
		DB:       nil,
		cache:    newCache(),
	}
	home, err := os.UserHomeDir()
	if err != nil {
		return &sqlite, err
	}
	path := filepath.Join(home, filename)
	db, err := sql.Open("sqlite", path)
	if err != nil {
		return &sqlite, err
	}
	sqlite.DB = db
	sqlite.filename = path
	return &sqlite, nil
}

// DefaultSqliteConn provides a connection for standard operation
func DefaultSqliteConn() (*Database, error) {
	return connect(DBName)
}

func testSqliteConn() (*Database, error) {
	return connect(DBNameTest)
}

// closeAndDelete will close a connection, and delete the sqlite file. used for testing
func (db *Database) closeAndDelete() error {
	err := db.DB.Close()
	if err != nil {
		return err
	}
	return os.Remove(db.filename)
}

// SaveChange will, given a title and seconds, create or update records accordingly
func (db *Database) SaveChange(windowTitle string, seconds int) error {
	window := NewWindow(windowTitle)

	err := window.safeWrite(db)
	if err != nil {
		return err
	}

	day := NewDay()
	err = day.safeWrite(db)
	if err != nil {
		return err
	}

	dw := NewDayWindow(day, window)
	inserted, err := dw.safeWrite(db)
	if err != nil {
		return err
	}
	if !inserted {
		return dw.AddSeconds(db, seconds)
	}

	return nil
}

// Init creates the database if it does not exist
func (db *Database) Init() error {
	initQuery := `
	CREATE TABLE IF NOT EXISTS window (
		id          INTEGER PRIMARY KEY AUTOINCREMENT,
		name        TEXT NOT NULL,
		UNIQUE      ( name )
	);
	CREATE TABLE IF NOT EXISTS day (
		id          INTEGER PRIMARY KEY AUTOINCREMENT,
		value       DATE NOT NULL,
		UNIQUE      ( value )
	);
	CREATE TABLE IF NOT EXISTS day_window (
		seconds     INT NOT NULL,
		day_id      INT NOT NULL,
		window_id   INT NOT NULL,
		PRIMARY KEY ( day_id, window_id ),
		FOREIGN KEY ( day_id )    REFERENCES day ( id ),
		FOREIGN KEY ( window_id ) REFERENCES window ( id )
	);
		`
	_, err := db.DB.Exec(initQuery)
	if err != nil {
		return err
	}
	return nil
}
