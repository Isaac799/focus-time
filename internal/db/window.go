package db

// Window is a database record of a desktop window
type Window struct {
	ID   int
	Name string
}

// NewWindow provides a window record for use in the database
func NewWindow(name string) Window {
	return Window{
		Name: cleanString(name),
	}
}

func (w *Window) read(db *Database) error {
	if len(w.Name) == 0 {
		return ErrWindowNameRequired
	}

	restored := db.cache.RestoreWindow(w)
	if restored {
		return nil
	}

	queryR := `
	SELECT id 
	FROM window
	WHERE name = $1
	`
	row := db.DB.QueryRow(queryR, w.Name)

	err := row.Scan(&w.ID)
	if err != nil {
		return err
	}

	db.cache.AddWindow(w)
	return nil
}

func (w *Window) safeWrite(db *Database) error {
	if len(w.Name) == 0 {
		return ErrWindowNameRequired
	}

	// reading will populate ID if exists
	// makes up for not re-inserting and returning id
	w.read(db)
	if w.ID != 0 {
		return nil
	}

	queryW := `
	INSERT INTO window (name) 
	VALUES ($1) 
	RETURNING id
	`
	row := db.DB.QueryRow(queryW, w.Name)
	err := row.Scan(&w.ID)
	if err != nil {
		return err
	}
	return nil
}
