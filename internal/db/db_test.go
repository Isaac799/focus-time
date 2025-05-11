package db

import (
	"fmt"
	"testing"

	_ "modernc.org/sqlite"
)

func TestDay_WriteAndRead(t *testing.T) {
	db, err := testSqliteConn()
	if err != nil {
		t.Fatal(err)
	}
	defer db.closeAndDelete()

	db.Init()

	day := NewDay()

	err = day.safeWrite(db)
	if err != nil {
		t.Fatal(err)
	}

	// handling conflict
	day2 := NewDay()

	err = day2.safeWrite(db)
	if err != nil {
		t.Fatal(err)
	}

	err = day.read(db)
	if err != nil {
		t.Fatal(err)
	}
	if day.ID == 0 {
		t.Fatal(ErrZeroID)
	}
	if day.Value.IsZero() {
		t.Fatal(ErrMissingValue)
	}
}

func TestWindow_WriteAndRead(t *testing.T) {
	db, err := testSqliteConn()
	if err != nil {
		t.Fatal(err)
	}
	defer db.closeAndDelete()

	db.Init()

	window := NewWindow("test window")

	err = window.safeWrite(db)
	if err != nil {
		db.closeAndDelete()
		t.Fatal(err)
	}

	// handling conflict
	window2 := NewWindow("test window")

	err = window2.safeWrite(db)
	if err != nil {
		db.closeAndDelete()
		t.Fatal(err)
	}

	err = window.read(db)
	if err != nil {
		db.closeAndDelete()
		t.Fatal(err)
	}
	if window.ID == 0 {
		db.closeAndDelete()
		t.Fatal(ErrZeroID)
	}
}

func setupTestDayWindow(t *testing.T, db *Database) *DayWindow {
	// Write a window for testing the assoc entity
	window := NewWindow("test window")

	err := window.safeWrite(db)
	if err != nil {
		db.closeAndDelete()
		t.Fatal(err)
	}

	// Write a day for testing the assoc entity
	day := NewDay()
	err = day.safeWrite(db)
	if err != nil {
		db.closeAndDelete()
		t.Fatal(err)
	}

	// Write assoc entity
	dw := NewDayWindow(day, window)

	return &dw
}

func TestDayWindow_WriteAndRead(t *testing.T) {
	db, err := testSqliteConn()
	if err != nil {
		db.closeAndDelete()
		t.Fatal(err)
	}
	defer db.closeAndDelete()

	db.Init()

	dw := setupTestDayWindow(t, db)

	_, err = dw.safeWrite(db)
	if err != nil {
		db.closeAndDelete()
		t.Fatal(err)
	}

	err = dw.read(db)
	if err != nil {
		db.closeAndDelete()
		t.Fatal(err)
	}
}

func TestDayWindow_AddSeconds(t *testing.T) {
	db, err := testSqliteConn()
	if err != nil {
		t.Fatal(err)
	}
	defer db.closeAndDelete()
	db.Init()

	dw := setupTestDayWindow(t, db)
	_, err = dw.safeWrite(db)
	if err != nil {
		db.closeAndDelete()
		t.Fatal(err)
	}

	// handling conflict
	dw2 := setupTestDayWindow(t, db)
	_, err = dw2.safeWrite(db)
	if err != nil {
		db.closeAndDelete()
		t.Fatal(err)
	}

	err = dw.AddSeconds(db, 200)
	if err != nil {
		db.closeAndDelete()
		t.Fatal(err)
	}

	if dw.Seconds != 200 {
		db.closeAndDelete()
		fmt.Println(dw.Seconds)
		t.Fatal(ErrUnexpectedValue)
	}

	err = dw.AddSeconds(db, 300)
	if err != nil {
		db.closeAndDelete()
		t.Fatal(err)
	}

	if dw.Seconds != 500 {
		db.closeAndDelete()
		t.Fatal(ErrUnexpectedValue)
	}
}
