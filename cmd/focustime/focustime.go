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
	database, err := db.DefaultSqliteConn()
	if err != nil {
		log.Fatal(err)
	}
	defer database.DB.Close()
	database.Init()

	action := flag.String("action", "watch", "watch | csv | print")

	// watch flags
	verbose := flag.Bool("verbose", false, "used with the 'watch' action to print current focus")

	// reporting defaults (last month, at least 1 minute long)
	now := time.Now()
	month := 30 * 24 * time.Hour
	ago := now.Add(-month)

	// reporting flags
	startStr := flag.String("start", ago.Format("2006-01-02"), "used in reporting to specify a start date")
	endStr := flag.String("end", now.Format("2006-01-02"), "used in reporting to specify an end date")
	duration := flag.Duration("duration", 1*time.Minute, "used in reporting to specify a minimum duration")

	flag.Parse()

	filter := db.NewReportFilter(startStr, endStr, duration)

	switch *action {
	case "watch":
		run(database, *verbose)
	case "print":
		database.PrintGroupedReport(filter)
	case "csv":
		database.WriteCSV(filter)
	default:
		log.Fatal("unknown action. try --help")
	}

	os.Exit(0)
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
