package log

import (
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"os"
	"sync"
)

type WebappLogger struct {
	StdoutLogger *slog.Logger
	StderrLogger *slog.Logger
	instanceId   string
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

	resp, err := http.Get("https://169.254.169.254/latest/meta-data/instance-id")
	globalLogger = &WebappLogger{}

	if err != nil || resp.StatusCode != http.StatusOK {
		globalLogger.instanceId = "localhost"
	} else {

		bodyReader := resp.Body
		defer bodyReader.Close()

		body, err := io.ReadAll(bodyReader)
		if err != nil {
			globalLogger.instanceId = "localhost"
		} else {
			globalLogger.instanceId = string(body)
		}
	}

	globalLogger.StdoutLogger = slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{AddSource: true}))
	globalLogger.StderrLogger = slog.New(slog.NewJSONHandler(os.Stderr, &slog.HandlerOptions{AddSource: true}))
	globalLogger.Info(fmt.Sprintf("Instance id is set to: %s", globalLogger.instanceId))
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
