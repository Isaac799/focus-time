// Package main is the primary entry point for the cli application
package main

import (
	"fmt"
	"log"
	"os"

	"github.com/Isaac799/focus-time/internal/focusedwindow"
	"github.com/Isaac799/focus-time/internal/prompt"
	"github.com/Isaac799/focus-time/internal/sqlite"
	_ "modernc.org/sqlite"
)

const (
	// Exit exits the program
	Exit = iota + 1
	// Upsert adds 1 second to tracked time
	Upsert
	// Read triggers a read from the sqlite database
	Read
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
		"Upsert",
		"Read",
	}

	sqlite.New()

	for {
		i := prompt.Select(options...)
		switch i {
		case Exit:
			os.Exit(0)
		case Upsert:
			sqlite.Upsert(title, 1)
		case Read:
			records, err := sqlite.Read()
			if err != nil {
				log.Fatal(err)
				return
			}
			for _, r := range records {
				fmt.Printf("%d: %s", r.Seconds, r.Title)
			}
		}
	}
}
