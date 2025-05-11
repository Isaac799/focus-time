// Package db keeps persistence
// syntax: https://github.com/mattn/go-sqlite3/blob/master/_example/simple/simple.go
// driver: https://pkg.go.dev/modernc.org/sqlite
package db

import "strings"

func cleanString(s string) string {
	remove := []string{"● "}
	for _, e := range remove {
		s = strings.ReplaceAll(s, e, "")
	}
	s = strings.ReplaceAll(s, "—", "-")
	s = strings.ReplaceAll(s, "  ", " ")
	s = strings.Trim(s, "\n\r\t ")
	return s
}
