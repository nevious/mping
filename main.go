package main

import (
	"fmt"
	"log/slog"
	"github.com/nevious/goping/internal/table"
)

func main() {
	records := []string{
		"1.1.1.1", "1.1.1.2", "185.178.192.107",
		"192.168.50.1", "192.168.50.22", "2a00:1450:400a:802::200e",
		"google.com", "nevious.ch",
	}

	if _, err := table.LaunchTablePing(table.MakeTable(records)); err != nil {
		slog.Error(
			"Unable to run programm",
			slog.String("Error:", fmt.Sprintf("%v", err)),
		)
	}
}
