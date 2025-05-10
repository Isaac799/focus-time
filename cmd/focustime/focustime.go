// Package main is the primary entry point for the cli application
package main

import (
	"fmt"
	"os"

	"github.com/Isaac799/focus-time/internal/focusedwindow"
	"github.com/Isaac799/focus-time/internal/prompt"
)

const (
	// Exit exits the program
	Exit = iota + 1
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

	options := []string{
		"Exit",
	}
	for {
		i := prompt.Select(options...)
		switch i {
		case Exit:
			os.Exit(0)
		}
	}
}
