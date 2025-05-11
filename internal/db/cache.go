package db

import "sync"

type cache struct {
	// protects below
	mu      sync.Mutex
	days    map[string]Day
	windows map[string]Window
}

func newCache() cache {
	return cache{
		days:    map[string]Day{},
		windows: map[string]Window{},
	}
}

// day

func (c *cache) AddDay(d *Day) {
	key := d.valueStr()

	if len(key) == 0 {
		return
	}

	c.mu.Lock()
	defer c.mu.Unlock()

	c.days[key] = *d
}

func (c *cache) ExistsDay(d *Day) bool {
	key := d.valueStr()

	if len(key) == 0 {
		return false
	}

	c.mu.Lock()
	defer c.mu.Unlock()

	_, exists := c.days[key]
	return exists
}

func (c *cache) RestoreDay(d *Day) bool {
	key := d.valueStr()

	if len(key) == 0 {
		return false
	}

	c.mu.Lock()
	defer c.mu.Unlock()

	day, exists := c.days[key]
	if !exists {
		return false
	}

	d.ID = day.ID
	return true
}

// window

func (c *cache) AddWindow(w *Window) {
	key := w.Name

	if len(key) == 0 {
		return
	}

	c.mu.Lock()
	defer c.mu.Unlock()

	c.windows[key] = *w
}

func (c *cache) ExistsWindow(w *Window) bool {
	key := w.Name

	if len(key) == 0 {
		return false
	}

	c.mu.Lock()
	defer c.mu.Unlock()

	_, exists := c.windows[key]
	return exists
}

func (c *cache) RestoreWindow(w *Window) bool {
	key := w.Name

	if len(key) == 0 {
		return false
	}

	c.mu.Lock()
	defer c.mu.Unlock()

	window, exists := c.windows[key]
	if !exists {
		return false
	}

	w.ID = window.ID
	return true
}
