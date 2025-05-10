// Package focusedwindow is a windows implementation of determining the currently
// focused window using user32 and syscall
package focusedwindow

import "errors"

var (
	// ErrNoFocusedWindow returned if no title is found for a focused window
	ErrNoFocusedWindow = errors.New("no focused window")
	// ErrNoTitle returned if a title is not found
	ErrNoTitle = errors.New("no title")
)
