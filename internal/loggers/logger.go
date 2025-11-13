package loggers

import (
	"encoding/json"
	"log"
	"qualifire-home-assignment/internal/models"

	"github.com/fatih/color"
)

type LoggerContract interface {
	Info()
	Error()
}

type Logger struct {
	Entry *models.LogEntry
}

func (l Logger) log(levelName string, colorFunc func(format string, a ...interface{})) {
	b, err := json.MarshalIndent(l.Entry, "", " ")
	if err != nil {
		log.Printf("logger marshal error: %v", err)
		return
	}
	colorFunc("%s: %s", levelName, string(b))
}

func (l Logger) Info() {
	colorFunc := color.New(color.FgGreen).PrintfFunc()
	l.log("INFO", colorFunc)
}

func (l Logger) Error() {
	colorFunc := color.New(color.FgRed).PrintfFunc()
	l.log("ERROR", colorFunc)
}
