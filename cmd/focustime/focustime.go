// Package main is the primary entry point for the cli application
package main

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/Isaac799/focus-time/internal/prompt"
	"github.com/Isaac799/focus-time/internal/sqlite"
	"github.com/Isaac799/focus-time/internal/watcher"
	_ "modernc.org/sqlite"
)

const (
	// Exit exits the program
	Exit = iota + 1
	// SeeCurrentFocus shows the currently focused wino name
	SeeCurrentFocus
)

func main() {
	options := []string{
		"Exit",
		"See Current Focus",
	}

	w := watcher.New()
	go consumeErr(&w)
	go consumeEvent(&w)
	go start(&w)

	for {
		i := prompt.Select(options...)
		switch i {
		case Exit:
			os.Exit(0)
		case SeeCurrentFocus:
			event := w.Read()
			fmt.Printf("Seconds: %d, Title: %s\n", int(event.Duration.Seconds()), event.Title)
		}
	}
}

func start(w *watcher.Watcher) {
	w.Start(time.Second * 1)
}

func consumeErr(w *watcher.Watcher) {
	for event := range w.OnError {
		fmt.Printf("Monitor Error Event: %s\n", event.Error())
	}
}

func consumeEvent(w *watcher.Watcher) {
	c, err := sqlite.DefaultSqliteConn()
	if err != nil {
		log.Fatal(err)
	}
	defer c.DB.Close()

	c.Init()

	for event := range w.OnChange {
		if event.Kind == watcher.FocusKindStart {
			fmt.Printf("Monitor Start: Seconds: %d, Title: %s\n", int(event.Duration.Seconds()), event.Title)
			continue
		}

		fmt.Printf("Monitor End: Seconds: %d, Title: %s\n", int(event.Duration.Seconds()), event.Title)
		err := c.SaveChange(event.Title, int(event.Duration.Seconds()))
		if err != nil {
			fmt.Printf("Monitor End Error: %s\n", err.Error())
		}
	}
}
