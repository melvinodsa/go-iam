package config

import (
	"os"

	"github.com/gofiber/fiber/v2/log"
)

type Logger struct {
	Level log.Level
}

func NewLogger(level log.Level) *Logger {
	log.SetLevel(level)
	log.SetOutput(os.Stdout)
	lgl := &Logger{
		Level: level,
	}
	return lgl
}
