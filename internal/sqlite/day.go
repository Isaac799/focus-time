package sqlite

import (
	"database/sql"
	"errors"
	"time"
)

// Day is a database record of a particular day
type Day struct {
	ID         int
	Value      time.Time
	InsertedAt time.Time
}

// NewDay provides a day record for use in the database
func NewDay() Day {
	return Day{
		Value: time.Now(),
	}
}

func (d *Day) read(c *Connection) error {
	queryR := `
SELECT id, inserted_at 
FROM day
WHERE value = $1
`
	row := c.db.QueryRow(queryR, d.Value)

	err := row.Scan(&d.ID, &d.InsertedAt)
	if err != nil {
		return err
	}
	return nil
}

func (d *Day) write(c *Connection) error {
	queryW := `
INSERT INTO day (value) 
VALUES ($1) 
RETURNING id
`
	row := c.db.QueryRow(queryW, d.Value)
	err := row.Scan(&d.ID)
	if err != nil {
		return err
	}
	return nil
}

// Save will save a day in the database, if it does not already exist
func (d *Day) Save(c *Connection) error {
	err := d.read(c)
	if err == nil {
		return nil
	}
	if !errors.Is(err, sql.ErrNoRows) {
		return err
	}
	return d.write(c)
}
