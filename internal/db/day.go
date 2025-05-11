package db

import (
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
	restored := db.cache.RestoreDay(d)
	if restored {
		return nil
	}

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

	db.cache.AddDay(d)
	return nil
}

func (d *Day) safeWrite(db *Database) error {
	// reading will populate ID if exists
	// makes up for not re-inserting and returning id
	d.read(db)
	if d.ID != 0 {
		return nil
	}

	queryW := `
	INSERT INTO day (value) 
	VALUES ($1) 
	RETURNING id
	`
	row := db.DB.QueryRow(queryW, d.valueStr())
	err := row.Scan(&d.ID)
	if err != nil {
		return err
	}
	return nil
}
