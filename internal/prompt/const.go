// Package prompt is a basic collection of cli interactions
package prompt

import "errors"

var (
	// ErrInvalidSelection is used if a prompt selection is out of range
	ErrInvalidSelection = errors.New("invalid selection")
)
