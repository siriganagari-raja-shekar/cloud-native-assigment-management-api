package log

import (
	"context"
	"errors"
	"log/slog"
	"os"
	"path/filepath"
	"runtime"
	"strconv"
	"sync"
	"time"
)

type WebappLogger struct {
	StdoutLogger *slog.Logger
	StderrLogger *slog.Logger
	leveler      *slog.LevelVar
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

	replace := func(groups []string, a slog.Attr) slog.Attr {

		// Remove the directory from the source's filename.
		if a.Key == slog.SourceKey {
			source := a.Value.Any().(*slog.Source)
			source.File = filepath.Base(source.File)
		}
		return a
	}

	lvl := &slog.LevelVar{}

	logLevelString := os.Getenv("LOG_LEVEL")
	logLevel, err := strconv.Atoi(logLevelString)
	if err == nil {
		if logLevel <= int(slog.LevelDebug) {
			lvl.Set(slog.LevelDebug)
		} else if logLevel <= int(slog.LevelInfo) {
			lvl.Set(slog.LevelInfo)
		} else if logLevel <= int(slog.LevelWarn) {
			lvl.Set(slog.LevelWarn)
		} else {
			lvl.Set(slog.LevelError)
		}
	} else {
		lvl.Set(slog.LevelWarn)
	}

	globalLogger = &WebappLogger{}
	globalLogger.leveler = lvl
	globalLogger.StdoutLogger = slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{AddSource: true, ReplaceAttr: replace}))
	globalLogger.StderrLogger = slog.New(slog.NewJSONHandler(os.Stderr, &slog.HandlerOptions{AddSource: true, ReplaceAttr: replace}))
}

func (logger *WebappLogger) Info(msg string) {
	_ = logger.log(msg, slog.LevelInfo)
}

func (logger *WebappLogger) Debug(msg string) {
	_ = logger.log(msg, slog.LevelDebug)
}

func (logger *WebappLogger) Warn(msg string) {
	_ = logger.log(msg, slog.LevelWarn)
}

func (logger *WebappLogger) Error(msg string) {
	_ = logger.log(msg, slog.LevelError)
}

func (logger *WebappLogger) log(msg string, level slog.Level) error {
	if level < logger.leveler.Level() {
		return errors.New("log level is higher")
	}
	var pcs [1]uintptr
	runtime.Callers(3, pcs[:]) // skip [Callers, Log<T>, getUpdatedRecord]
	r := slog.NewRecord(time.Now(), level, msg, pcs[0])
	var err error
	if level <= slog.LevelInfo {
		err = logger.StdoutLogger.Handler().Handle(context.Background(), r)
	} else {
		err = logger.StderrLogger.Handler().Handle(context.Background(), r)
	}
	return err
}
