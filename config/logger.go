package config

import (
	"os"

	"github.com/gofiber/fiber/v2/log"
)

// Logger holds logging configuration settings.
type Logger struct {
	Level log.Level // The minimum log level to output
}

// NewLogger creates a new Logger instance with the specified log level.
// It configures the global logger to output to stdout with the given level.
//
// Parameters:
//   - level: The minimum log level to output (e.g., log.LevelInfo, log.LevelDebug)
//
// Returns a configured Logger instance.
func NewLogger(level log.Level) *Logger {
	log.SetLevel(level)
	log.SetOutput(os.Stdout)
	lgl := &Logger{
		Level: level,
	}
	return lgl
}
