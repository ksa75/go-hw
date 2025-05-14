package logger

import (
	"io"
	"log"
	"os"
	"strings"
)

type Level int

const (
	DEBUG Level = iota
	INFO
	ERROR
)

type Logger struct {
	level  Level
	logger *log.Logger
}

// New создаёт логгер с заданным уровнем и файлом.
func New(levelStr, logPath string) (*Logger, error) {
	level := parseLevel(levelStr)

	// открываем файл
	logFile, err := os.OpenFile(logPath, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0o644)
	if err != nil {
		return nil, err
	}

	multiWriter := io.MultiWriter(os.Stdout, logFile)

	return &Logger{
		level:  level,
		logger: log.New(multiWriter, "", 0),
	}, nil
}

func (l *Logger) Printf(format string, v ...any) {
	if l.level <= INFO {
		l.logger.Printf(format, v...)
	}
}

func (l *Logger) Debug(msg string) {
	if l.level <= DEBUG {
		l.logger.Printf("[DEBUG] %s", msg)
	}
}

func (l *Logger) Info(msg string) {
	if l.level <= INFO {
		l.logger.Printf("[INFO] %s", msg)
	}
}

func (l *Logger) Error(msg string) {
	if l.level <= ERROR {
		l.logger.Printf("[ERROR] %s", msg)
	}
}

func parseLevel(s string) Level {
	switch strings.ToLower(s) {
	case "debug":
		return DEBUG
	case "error":
		return ERROR
	default:
		return INFO
	}
}
