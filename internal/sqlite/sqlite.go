package sqlite

import (
	"database/sql"
	"os"
	"path/filepath"
)

// Connection stores the database connection and name
type Connection struct {
	filename string
	DB       *sql.DB
}

func connect(filename string) (*Connection, error) {
	sqlite := Connection{
		filename: "",
		DB:       nil,
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
func DefaultSqliteConn() (*Connection, error) {
	return connect(DBName)
}

func testSqliteConn() (*Connection, error) {
	return connect(DBNameTest)
}

// closeAndDelete will close a connection, and delete the sqlite file. used for testing
func (c *Connection) closeAndDelete() error {
	err := c.DB.Close()
	if err != nil {
		return err
	}
	return os.Remove(c.filename)
}

// SaveChange will, given a title and seconds, create or update records accordingly
func (c *Connection) SaveChange(windowTitle string, seconds int) error {
	window := NewWindow(windowTitle)

	err := window.safeWrite(c)
	if err != nil {
		return err
	}

	day := NewDay()
	err = day.safeWrite(c)
	if err != nil {
		return err
	}

	dw := NewDayWindow(day, window)
	inserted, err := dw.safeWrite(c)
	if err != nil {
		return err
	}
	if !inserted {
		return dw.AddSeconds(c, seconds)
	}

	return nil
}

// Init creates the database if it does not exist
func (c *Connection) Init() error {
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
	_, err := c.DB.Exec(initQuery)
	if err != nil {
		return err
	}
	return nil
}
