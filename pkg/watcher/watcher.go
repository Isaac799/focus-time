// Package watcher enables watching the focused window
package watcher

import (
	"errors"
	"time"

	"github.com/Isaac799/focus-time/pkg/focusedwindow"
)

const (
	// FocusKindStart is sent via chan when a window focus starts
	FocusKindStart = iota
	// FocusKindEnd is sent via chan when a window focus ends
	FocusKindEnd
)

// Watcher will watch the focused window for changes
type Watcher struct {
	alive       bool
	windowTitle string
	focusStart  time.Time
	OnChange    chan Event
	OnError     chan error
}

// New provides a monitor to watch the focused window for changes
func New() Watcher {
	return Watcher{
		alive:      true,
		focusStart: time.Now(),
		OnChange:   make(chan Event),
		OnError:    make(chan error),
	}
}

// Event is when the active window changes
type Event struct {
	Kind     int
	Title    string
	Duration time.Duration
}

// Start spins up a routine to check the window title and update
func (w *Watcher) Start(frequency time.Duration) {
	focusedWindow := focusedwindow.New()

	for {
		if !w.alive {
			break
		}
		time.Sleep(frequency)

		title, err := focusedWindow.Title()

		focusToNone := err != nil && len(w.windowTitle) > 0 && errors.Is(err, focusedwindow.ErrNoFocusedWindow)
		if focusToNone {
			w.change(title)
			continue
		}

		if err != nil {
			if len(w.windowTitle) == 0 {
				continue
			}
			w.OnError <- err
			continue
		}

		starting := len(title) > 0 && len(w.windowTitle) == 0
		if starting {
			w.windowTitle = title
			w.OnChange <- Event{
				Kind:     FocusKindStart,
				Duration: time.Duration(0),
				Title:    title,
			}
			continue
		}

		if w.windowTitle == title {
			continue
		}
		w.change(title)
		continue
	}

	close(w.OnChange)
	close(w.OnChange)
}

func (w *Watcher) change(toTitle string) {
	w.OnChange <- Event{
		Kind:     FocusKindEnd,
		Duration: time.Since(w.focusStart),
		Title:    w.windowTitle,
	}
	w.windowTitle = toTitle
	w.focusStart = time.Now()

}

// Stop will end the watcher loop and close its channels
func (w *Watcher) Stop() {
	w.alive = false
}

func (w *Watcher) Read() Event {
	return Event{
		Kind:     FocusKindEnd,
		Duration: time.Since(w.focusStart),
		Title:    w.windowTitle,
	}
}
