// Package sqlite keeps persistence
// syntax: https://github.com/mattn/go-sqlite3/blob/master/_example/simple/simple.go
// driver: https://pkg.go.dev/modernc.org/sqlite
package sqlite

import (
	"database/sql"
	"fmt"
	"os"
	"path/filepath"
	"time"
)

// Record is a database record
type Record struct {
	id      uint
	Title   string
	Seconds uint
	iat     time.Time
}

func db() (*sql.DB, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return nil, err
	}
	path := filepath.Join(home, ".focustime.db")
	db, err := sql.Open("sqlite", path)
	if err != nil {
		return nil, err
	}
	return db, nil
}

// New ensures a new database
func New() error {
	db, err := db()
	if err != nil {
		return err
	}
	defer db.Close()

	q := `
	CREATE TABLE IF NOT EXISTS focustime (
		id INTEGER PRIMARY KEY AUTOINCREMENT, 
		title TEXT, 
		seconds INT, 
		iat DATETIME DEFAULT CURRENT_TIMESTAMP 
	);
		`
	_, err = db.Exec(q)
	if err != nil {
		return err
	}
	return nil
}

func find(title string) (Record, error) {
	record := Record{}
	db, err := db()
	if err != nil {
		return record, err
	}
	defer db.Close()

	q := `
	SELECT id, seconds, iat FROM focustime WHERE title = $1;
	`
	row := db.QueryRow(q, title)
	if row == nil {
		return record, nil
	}
	row.Scan(&record.id, &record.Seconds, &record.iat)
	return record, nil
}

// Upsert will insert a record,to track time, or update an existing one
func Upsert(title string, seconds uint) error {
	record, err := find(title)
	if err != nil {
		return err
	}

	db, err := db()
	if err != nil {
		return err
	}
	defer db.Close()

	if record.id == 0 {
		fmt.Println("creating")
		q := `
		INSERT INTO focustime ( title, seconds ) VALUES ( $1, $2 );
		`
		_, err := db.Exec(q, title, seconds)
		if err != nil {
			return err
		}
		return nil
	}

	fmt.Println("updating")
	q := `
	UPDATE focustime SET seconds = $2 WHERE id = $1;
	`
	_, err = db.Exec(q, record.id, record.Seconds+seconds)
	if err != nil {
		return err
	}

	return nil
}

// Read gives a list of records
func Read() ([]Record, error) {
	db, err := db()
	if err != nil {
		return nil, err
	}
	defer db.Close()

	q := `
	SELECT id, title, seconds, iat FROM focustime
	`
	rows, err := db.Query(q)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	records := []Record{}
	for rows.Next() {
		record := Record{}
		err := rows.Scan(&record.id, &record.Title, &record.Seconds, &record.iat)
		if err != nil {
			return nil, err
		}
		records = append(records, record)
	}
	return records, nil
}
