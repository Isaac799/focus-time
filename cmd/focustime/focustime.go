// Package main is the primary entry point for the cli application
package main

import (
	"fmt"
	"os"

	"github.com/Isaac799/focus-time/internal/focusedwindow"
	"github.com/Isaac799/focus-time/internal/prompt"
	"github.com/Isaac799/focus-time/internal/sqlite"
	_ "modernc.org/sqlite"
)

const (
	// Exit exits the program
	Exit = iota + 1
)

func main() {
	w := focusedwindow.New()
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

	sqlite.Init()

	for {
		i := prompt.Select(options...)
		switch i {
		case Exit:
			os.Exit(0)
		}
	}
}
