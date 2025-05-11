package sqlite

import (
	"strings"
	"time"
)

// ReportItem is an item in a report
type ReportItem struct {
	Title   string
	When    time.Time
	Seconds int
}

// Report is a report derived from joined tables
type Report struct {
	Items []ReportItem
}

// Report provides all records with at least certain amount of duration time tracked
func (c *Connection) Report(duration time.Duration) (*Report, error) {
	queryAll := `
	SELECT w.name, d.value, dw.seconds
	FROM day_window dw
	LEFT JOIN window w ON w.id = dw.window_id 
	LEFT JOIN day d ON d.id = dw.day_id 
	WHERE dw.seconds > $1
	`
	rows, err := c.DB.Query(queryAll, duration.Seconds())
	if err != nil {
		return nil, err
	}

	items := []ReportItem{}
	for rows.Next() {
		item := ReportItem{}
		err := rows.Scan(&item.Title, &item.When, &item.Seconds)
		if err != nil {
			return nil, err
		}
		items = append(items, item)
	}
	report := Report{items}
	return &report, nil
}

}
