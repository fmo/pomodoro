package main

import (
	"log/slog"
	"os"

	"github.com/fmo/pomodoro/cmd"
)

func main() {
	handler := slog.NewJSONHandler(os.Stderr, &slog.HandlerOptions{
		Level: slog.LevelDebug,
	})

	logger := slog.New(handler)

	slog.SetDefault(logger)

	cmd.Execute()
}
