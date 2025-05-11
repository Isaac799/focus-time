// Package main is the primary entry point for the cli application
package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/Isaac799/focus-time/internal/db"
	"github.com/Isaac799/focus-time/pkg/watcher"
	_ "modernc.org/sqlite"
)

func main() {
	db, err := db.DefaultSqliteConn()
	if err != nil {
		log.Fatal(err)
	}
	defer db.DB.Close()
	db.Init()

	action := flag.String("action", "watch", "watch | csv | print")
	verbose := flag.Bool("verbose", false, "used with the 'watch' action to print current focus")
	flag.Parse()

	switch *action {
	case "watch":
		run(db, *verbose)
	case "print":
		db.PrintGroupedReport()
	case "csv":
		db.WriteCSV()
	case "help":
		printHelp()
	default:
		printHelp()
		os.Exit(1)
	}

	os.Exit(0)
}

func printHelp() {
	fmt.Println("")
	fmt.Println("acceptable actions: watch | csv | print")
	fmt.Println("- watch   is to be kept open while using computer")
	fmt.Println("- csv   spits out a csv in working dir")
	fmt.Println("- print spits out a report in the console, useful with grep")

}

func run(db *db.Database, verbose bool) {
	w := watcher.New()
	go consumeErr(&w)
	go consumeEvent(&w, db)
	go start(&w)

	fmt.Println("watcher active")
	fmt.Println("keep window open")
	fmt.Println("press [ctrl + c] to quit")

	if verbose {
		fmt.Println("running in verbose mode")
	}

	for {
		if !verbose {
			time.Sleep(60 * time.Second)
			continue
		}

		time.Sleep(1 * time.Second)
		event := w.Read()
		fmt.Printf("%d %s\n", int(event.Duration.Seconds()), event.Title)
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

func consumeEvent(w *watcher.Watcher, db *db.Database) {
	for event := range w.OnChange {
		if event.Kind == watcher.FocusKindStart {
			continue
		}
		err := db.SaveChange(event.Title, int(event.Duration.Seconds()))
		if err != nil {
			fmt.Printf("watcher focus end error: %s\n", err.Error())
		}
	}
}
