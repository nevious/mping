package main

import (
	"fmt"
	"log/slog"
	"github.com/nevious/goping/internal/table"
	"github.com/nevious/goping/internal/parser"
)

func main() {
	records := parser.Parse()

	if _, err := table.LaunchTablePing(table.MakeTable(records)); err != nil {
		slog.Error(
			"Unable to run programm",
			slog.String("Error:", fmt.Sprintf("%v", err)),
		)
	}
}
