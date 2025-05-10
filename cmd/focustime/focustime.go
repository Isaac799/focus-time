// Package main is the primary entry point for the cli application
package main

import (
	"fmt"

	"github.com/Isaac799/focus-time/internal/focusedwindow"
)

func main() {
	w := focusedwindow.NewFocusedWindow()
	title, err := w.Title()
	if err != nil {
		fmt.Println("failed")
		fmt.Println(err.Error())
		return
	}
	fmt.Println("success")
	fmt.Println(title)
}
