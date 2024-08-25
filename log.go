package main

import (
	"log/slog"
	"os"

	"github.com/lmittmann/tint"
	"github.com/mattn/go-isatty"
)

func newLog(opts slog.HandlerOptions) (*slog.Logger, error) {
	w := os.Stderr
	logger := slog.New(
		tint.NewHandler(w, &tint.Options{
			Level:      opts.Level,
			AddSource:  opts.AddSource,
			NoColor:    !isatty.IsTerminal(w.Fd()),
			TimeFormat: "2006-01-02T15:04:05.000Z",
		}),
	)
	return logger, nil
}
