package main

import (
	"fmt"
	"log/slog"
	"runtime"
	"github.com/nevious/mping/internal/views"
	"github.com/nevious/mping/internal/parser"
)

func main() {
	records := parser.Parse()

	// unpriv ping socket don't work for me, even with set ping_group_range
	switch runtime.GOOS {
		case "darwin", "ios":
		case "linux":
			slog.Info("you may need to adjust the net.ipv4.ping_group_range kernel state")
	}

	if _, err := views.LaunchTablePing(views.MakeTable(records)); err != nil {
		slog.Error(
			"Unable to run programm",
			slog.String("Error:", fmt.Sprintf("%v", err)),
		)
	}
}
