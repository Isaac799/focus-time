package sqlite

import (
	"testing"

	_ "modernc.org/sqlite"
)

func TestDay_WriteAndRead(t *testing.T) {
	c, err := testSqliteConn()
	if err != nil {
		t.Fatal(err)
	}
	defer c.closeAndDelete()

	c.Init()

	day := NewDay()

	err = day.safeWrite(c)
	if err != nil {
		t.Fatal(err)
	}

	// handling conflict
	day2 := NewDay()

	err = day2.safeWrite(c)
	if err != nil {
		t.Fatal(err)
	}

	err = day.read(c)
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
	c, err := testSqliteConn()
	if err != nil {
		t.Fatal(err)
	}
	defer c.closeAndDelete()

	c.Init()

	window := NewWindow("test window")

	err = window.safeWrite(c)
	if err != nil {
		t.Fatal(err)
	}

	// handling conflict
	window2 := NewWindow("test window")

	err = window2.safeWrite(c)
	if err != nil {
		t.Fatal(err)
	}

	err = window.read(c)
	if err != nil {
		t.Fatal(err)
	}
	if window.ID == 0 {
		t.Fatal(ErrZeroID)
	}
}

func setupTestDayWindow(t *testing.T, c *Connection) *DayWindow {
	// Write a window for testing the assoc entity
	window := NewWindow("test window")

	err := window.safeWrite(c)
	if err != nil {
		t.Fatal(err)
	}

	// Write a day for testing the assoc entity
	day := NewDay()
	err = day.safeWrite(c)
	if err != nil {
		t.Fatal(err)
	}

	// Write assoc entity
	dw := NewDayWindow(day, window)

	return &dw
}

func TestDayWindow_WriteAndRead(t *testing.T) {
	c, err := testSqliteConn()
	if err != nil {
		t.Fatal(err)
	}
	defer c.closeAndDelete()

	c.Init()

	dw := setupTestDayWindow(t, c)

	_, err = dw.safeWrite(c)
	if err != nil {
		t.Fatal(err)
	}

	err = dw.read(c)
	if err != nil {
		t.Fatal(err)
	}
}

func TestDayWindow_AddSeconds(t *testing.T) {
	c, err := testSqliteConn()
	if err != nil {
		t.Fatal(err)
	}
	defer c.closeAndDelete()
	c.Init()

	dw := setupTestDayWindow(t, c)
	_, err = dw.safeWrite(c)
	if err != nil {
		t.Fatal(err)
	}

	// handling conflict
	dw2 := setupTestDayWindow(t, c)
	_, err = dw2.safeWrite(c)
	if err != nil {
		t.Fatal(err)
	}

	err = dw.AddSeconds(c, 200)
	if err != nil {
		t.Fatal(err)
	}

	if dw.Seconds != 200 {
		t.Fatal(ErrUnexpectedValue)
	}

	err = dw.AddSeconds(c, 300)
	if err != nil {
		t.Fatal(err)
	}

	if dw.Seconds != 500 {
		t.Fatal(ErrUnexpectedValue)
	}
}
