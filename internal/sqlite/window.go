package sqlite

import (
	"database/sql"
	"errors"
)

// Window is a database record of a desktop window
type Window struct {
	ID   int
	Name string
}

// NewWindow provides a window record for use in the database
func NewWindow(name string) Window {
	return Window{
		Name: name,
	}
}

func (w *Window) read(c *Connection) error {
	if len(w.Name) == 0 {
		return ErrWindowNameRequired
	}

	queryR := `
SELECT id 
FROM window
WHERE name = $1
`
	row := c.DB.QueryRow(queryR, w.Name)

	err := row.Scan(&w.ID)
	if err != nil {
		return err
	}
	return nil
}

func (w *Window) safeWrite(c *Connection) error {
	if len(w.Name) == 0 {
		return ErrWindowNameRequired
	}

	err := w.read(c)
	if !errors.Is(err, sql.ErrNoRows) {
		return err
	}

	queryW := `
INSERT INTO window (name) 
VALUES ($1) 
RETURNING id
`
	row := c.DB.QueryRow(queryW, w.Name)
	err = row.Scan(&w.ID)
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
	return w.safeWrite(c)
}
