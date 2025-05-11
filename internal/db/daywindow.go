package db

import (
	"database/sql"
	"errors"
)

// DayWindow is a database record
type DayWindow struct {
	Seconds  int
	DayID    int
	WindowID int
}

// NewDayWindow is a associative database record of a particular day and window
func NewDayWindow(day Day, window Window) DayWindow {
	return DayWindow{
		DayID:    day.ID,
		WindowID: window.ID,
	}
}

func (dw *DayWindow) read(db *Database) error {
	queryR := `
SELECT seconds 
FROM day_window
WHERE day_id = $1 AND window_id = $2
`
	row := db.DB.QueryRow(queryR, dw.DayID, dw.WindowID)

	err := row.Scan(&dw.Seconds)
	if err != nil {
		return err
	}
	return nil
}

func (dw *DayWindow) safeWrite(db *Database) (bool, error) {
	err := dw.read(db)
	if !errors.Is(err, sql.ErrNoRows) {
		return false, err
	}

	queryW := `
INSERT INTO day_window (day_id, window_id, seconds) 
VALUES ($1, $2, $3) 
`
	_, err = db.DB.Exec(queryW, dw.DayID, dw.WindowID, dw.Seconds)
	return true, err
}

// AddSeconds grabs current seconds from db, adds n, and updates the record
func (dw *DayWindow) AddSeconds(db *Database, n int) error {
	err := dw.read(db)
	if err != nil {
		return err
	}

	sum := dw.Seconds + n

	queryUpdate := `
UPDATE day_window SET seconds = $1
WHERE day_id = $2 AND window_id = $3
`
	_, err = db.DB.Exec(queryUpdate, sum, dw.DayID, dw.WindowID)
	if err != nil {
		return err
	}

	dw.Seconds = sum
	return nil
}
