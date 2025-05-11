package db

import (
	"database/sql"
	"errors"
	"time"
)

// Day is a database record of a particular day
type Day struct {
	ID    int
	Value time.Time
}

// NewDay provides a day record for use in the database
func NewDay() Day {
	return Day{
		Value: time.Now(),
	}
}

func (d *Day) valueStr() string {
	return d.Value.Format("2006-01-02")
}

func (d *Day) read(db *Database) error {
	queryR := `
SELECT id 
FROM day
WHERE value = $1
`
	row := db.DB.QueryRow(queryR, d.valueStr())

	err := row.Scan(&d.ID)
	if err != nil {
		return err
	}
	return nil
}

func (d *Day) safeWrite(db *Database) error {
	err := d.read(db)
	if !errors.Is(err, sql.ErrNoRows) {
		return err
	}

	queryW := `
INSERT INTO day (value) 
VALUES ($1) 
RETURNING id
`
	row := db.DB.QueryRow(queryW, d.valueStr())
	err = row.Scan(&d.ID)
	if err != nil {
		return err
	}
	return nil
}

// Save will save a day in the database, if it does not already exist
func (d *Day) Save(db *Database) error {
	err := d.read(db)
	if err == nil {
		return nil
	}
	if !errors.Is(err, sql.ErrNoRows) {
		return err
	}
	return d.safeWrite(db)
}
