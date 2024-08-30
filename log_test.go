package main

import (
	"log/slog"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestNew(t *testing.T) {
	opts := slog.HandlerOptions{
		Level:     slog.LevelDebug,
		AddSource: false,
	}
	log, err := newLog(opts)
	require.NoError(t, err)
	require.NotNil(t, log)
	log.Debug("debug message", "key", "value")
	log.Info("info message", "key", "value")
	log.Warn("warn message", "key", "value")
	log.Error("error message", "key", "value")
}
