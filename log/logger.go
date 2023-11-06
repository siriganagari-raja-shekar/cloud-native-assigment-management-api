package log

import (
	"log/slog"
	"os"
	"sync"
)

type WebappLogger struct {
	StdoutLogger *slog.Logger
	StderrLogger *slog.Logger
}

var (
	globalLogger *WebappLogger
	once         sync.Once
)

func GetLoggerInstance() *WebappLogger {
	once.Do(createLoggerInstance)
	return globalLogger
}

func createLoggerInstance() {
	globalLogger = &WebappLogger{
		StdoutLogger: slog.New(slog.NewJSONHandler(os.Stdout, nil)),
		StderrLogger: slog.New(slog.NewJSONHandler(os.Stderr, nil)),
	}
}

func (logger *WebappLogger) Info(msg string) {
	logger.StdoutLogger.Info(msg)
}

func (logger *WebappLogger) Debug(msg string) {
	logger.StdoutLogger.Debug(msg)
}

func (logger *WebappLogger) Warn(msg string) {
	logger.StderrLogger.Warn(msg)
}

func (logger *WebappLogger) Error(msg string) {
	logger.StderrLogger.Error(msg)
}
