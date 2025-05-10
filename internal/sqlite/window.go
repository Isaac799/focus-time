package sqlite

import (
	"database/sql"
	"errors"
	"time"
)

// Window is a database record of a desktop window
type Window struct {
	ID         int
	Name       string
	InsertedAt time.Time
}

// NewWindow provides a window record for use in the database
func NewWindow() Window {
	return Window{}
}

func (w *Window) read(c *Connection) error {
	if len(w.Name) == 0 {
		return ErrWindowNameRequired
	}

	queryR := `
SELECT id, inserted_at 
FROM window
WHERE name = $1
`
	row := c.db.QueryRow(queryR, w.Name)

	err := row.Scan(&w.ID, &w.InsertedAt)
	if err != nil {
		return err
	}
	return nil
}

func (w *Window) write(c *Connection) error {
	if len(w.Name) == 0 {
		return ErrWindowNameRequired
	}

	queryW := `
INSERT INTO window (name) 
VALUES ($1) 
RETURNING id
`
	row := c.db.QueryRow(queryW, w.Name)
	err := row.Scan(&w.ID)
	if err != nil {
		return err
	}
	return nil
}

// Save will save a window in the database, if it does not already exist
func (w *Window) Save(c *Connection) error {
	err := w.read(c)
	if err == nil {
		return nil
	}
	if !errors.Is(err, sql.ErrNoRows) {
		return err
	}
	return w.write(c)
}
