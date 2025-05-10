package sqlite

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

func (dw *DayWindow) read(c *Connection) error {
	queryR := `
SELECT seconds 
FROM day_window
WHERE day_id = $1 AND window_id = $2
`
	row := c.DB.QueryRow(queryR, dw.DayID, dw.WindowID)

	err := row.Scan(&dw.Seconds)
	if err != nil {
		return err
	}
	return nil
}

func (dw *DayWindow) write(c *Connection) error {
	queryW := `
INSERT INTO day_window (day_id, window_id, seconds) 
VALUES ($1, $2, $3) 
`
	_, err := c.DB.Exec(queryW, dw.DayID, dw.WindowID, dw.Seconds)
	return err
}

// AddSeconds grabs current seconds from db, adds n, and updates the record
func (dw *DayWindow) AddSeconds(c *Connection, n int) error {
	err := dw.read(c)
	if err != nil {
		return err
	}

	sum := dw.Seconds + n

	queryUpdate := `
UPDATE day_window SET seconds = $1
WHERE day_id = $2 AND window_id = $3
`
	_, err = c.DB.Exec(queryUpdate, sum, dw.DayID, dw.WindowID)
	if err != nil {
		return err
	}

	dw.Seconds = sum
	return nil
}
