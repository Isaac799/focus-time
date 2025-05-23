package db

import (
	"encoding/csv"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"text/tabwriter"
	"time"
)

// ReportItem is an item in a report
type ReportItem struct {
	Title   string
	When    time.Time
	Seconds int
}

// ReportFilter is the filters that can be applied to a report
type ReportFilter struct {
	start           time.Time
	end             time.Time
	atLeastDuration time.Duration
}

// NewReportFilter takes in flags and generates a report filter that can be used
func NewReportFilter(startDateStr, endDateStr *string, duration *time.Duration) ReportFilter {
	filter := ReportFilter{}
	startDate, err := time.Parse("2006-01-02", *startDateStr)
	if err == nil {
		filter.start = startDate
	}
	endDate, err := time.Parse("2006-01-02", *endDateStr)
	if err == nil {
		filter.end = endDate
	}
	filter.atLeastDuration = *duration
	return filter
}

// Report is a report derived from joined tables
type Report struct {
	Items []ReportItem
}

// Report provides all records with at least certain amount of duration time tracked
func (db *Database) Report(filter ReportFilter) (*Report, error) {
	queryAll := `
	SELECT w.name, d.value, dw.seconds
	FROM day_window dw
	LEFT JOIN window w ON w.id = dw.window_id 
	LEFT JOIN day d ON d.id = dw.day_id 
	WHERE dw.seconds > $1 
	AND   d.value >= $2
	AND   d.value <= $3
	`
	rows, err := db.DB.Query(
		queryAll,
		filter.atLeastDuration.Seconds(),
		filter.start.Format("2006-01-02"),
		filter.end.Format("2006-01-02"),
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

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

// GroupedByTitleSuffix will attempt to group items by their suffix according to common delimiters
func (r *Report) GroupedByTitleSuffix() map[string][]ReportItem {
	m := map[string][]ReportItem{}
	delimiters := []string{"-", "|"}

	for _, e := range r.Items {
		candidates := map[string]bool{}

		for _, deli := range delimiters {
			parts := strings.Split(e.Title, deli)
			if len(parts) == 1 {
				continue
			}
			suffix := parts[len(parts)-1]
			candidates[suffix] = true
		}

		key := ""
		for c := range candidates {
			if len(c) <= len(key) {
				continue
			}
			key = c
		}

		key = cleanString(key)
		if key == "" {
			key = "No Group"
		}

		if m[key] == nil {
			m[key] = []ReportItem{}
		}

		// Remove suffix from string
		e.Title = strings.Replace(e.Title, key, "", 1)
		e.Title = cleanString(e.Title)
		for _, deli := range delimiters {
			e.Title = strings.TrimSuffix(e.Title, deli)
		}

		m[key] = append(m[key], e)
	}
	return m
}

// PrintGroupedReport will print a report of focused windows, grouped by suffix
func (db *Database) PrintGroupedReport(filter ReportFilter) {
	records, err := db.Report(filter)
	if err != nil {
		fmt.Print(err)
		return
	}
	writer := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
	fmt.Fprintln(writer, "\nGroup\tTitle\tWhen\tDuration")
	for key, items := range records.GroupedByTitleSuffix() {
		s := fmt.Sprintf("%s\t-\t-\t-", key)
		fmt.Fprintln(writer, s)
		for _, e := range items {
			dur := time.Duration(e.Seconds) * time.Second
			s := fmt.Sprintf("%s\t%s\t%s\t%s", key, e.Title, e.When.Format("2006-01-02"), dur.String())
			fmt.Fprintln(writer, s)
		}
	}
	writer.Flush()
}

// WriteCSV will write a report of focused windows, grouped by suffix, to current working dir
func (db *Database) WriteCSV(filter ReportFilter) {
	wd, err := os.Getwd()
	if err != nil {
		fmt.Print(err)
		return
	}
	p := filepath.Join(wd, CSVFileName)
	file, err := os.Create(p)
	if err != nil {
		fmt.Print(err)
		return
	}
	defer file.Close()

	records, err := db.Report(filter)
	if err != nil {
		fmt.Print(err)
		return
	}

	writer := csv.NewWriter(file)
	defer writer.Flush()

	header := []string{"Group", "Title", "When", "Duration"}
	err = writer.Write(header)
	if err != nil {
		fmt.Print(err)
		return
	}

	for key, items := range records.GroupedByTitleSuffix() {
		for _, e := range items {
			row := make([]string, 4)
			dur := time.Duration(e.Seconds) * time.Second
			row[0] = key
			row[1] = e.Title
			row[2] = e.When.Format("2006-01-02")
			row[3] = dur.String()
			err = writer.Write(row)
			if err != nil {
				fmt.Print(err)
				return
			}
		}
	}
}
