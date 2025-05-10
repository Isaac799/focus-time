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

func (d *Day) valueStr() string {
	return d.Value.Local().Format("2006-01-02")
}

func (d *Day) read(c *Connection) error {
	queryR := `
SELECT id, inserted_at 
FROM day
WHERE value = $1
`
	row := c.DB.QueryRow(queryR, d.valueStr())

	err := row.Scan(&d.ID, &d.InsertedAt)
	if err != nil {
		return err
	}
	return nil
}

func (d *Day) safeWrite(c *Connection) error {
	err := d.read(c)
	if !errors.Is(err, sql.ErrNoRows) {
		return err
	}

	queryW := `
INSERT INTO day (value) 
VALUES ($1) 
RETURNING id
`
	row := c.DB.QueryRow(queryW, d.valueStr())
	err = row.Scan(&d.ID)
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
	return d.safeWrite(c)
}
