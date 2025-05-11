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
	// SeeReport shows a summary of the time tracked
	SeeReport
	// SeeReportGrouped shows a summary of the time tracked, attempting to be grouped by title
	SeeReportGrouped
)

func main() {
	options := []string{
		"Exit",
		"See Current Focus",
		"Report",
		"Report Grouped",
	}

	c, err := sqlite.DefaultSqliteConn()
	if err != nil {
		log.Fatal(err)
	}
	defer c.DB.Close()
	c.Init()

	w := watcher.New()
	go consumeErr(&w)
	go consumeEvent(&w, c)
	go start(&w)

	for {
		i := prompt.Select(options...)
		switch i {
		case Exit:
			os.Exit(0)
		case SeeCurrentFocus:
			event := w.Read()
			fmt.Printf("Seconds: %d, Title: %s\n", int(event.Duration.Seconds()), event.Title)
		case SeeReport:
			c.PrintReport()
		case SeeReportGrouped:
			c.PrintGroupedReport()
		}
	}
}

func start(w *watcher.Watcher) {
	w.Start(time.Second * 1)
}

func consumeErr(w *watcher.Watcher) {
	for event := range w.OnError {
		fmt.Printf("watcher error: %s\n", event.Error())
	}
}

func consumeEvent(w *watcher.Watcher, c *sqlite.Connection) {
	for event := range w.OnChange {
		if event.Kind == watcher.FocusKindStart {
			continue
		}
		err := c.SaveChange(event.Title, int(event.Duration.Seconds()))
		if err != nil {
			fmt.Printf("watcher focus end error: %s\n", err.Error())
		}
	}
}
