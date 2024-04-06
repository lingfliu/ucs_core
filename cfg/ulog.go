package cfg

import (
	"context"
	"io"
	"log/slog"
	"os"
)

const LOG_LEVEL_INFO = 1
const LOG_LEVEL_DEBUG = 2
const LOG_LEVEL_WARN = 3
const LOG_LEVEL_ERROR = 4

type ULogger struct {
	Tag     string
	_logger *slog.Logger
}

var _ulogger *ULogger = nil

func GetULogger() *ULogger {
	if _ulogger == nil {
		_ulogger = &ULogger{}
	}
	return _ulogger
}

func (l *ULogger) Config(level int, exportPath string, hasTimestamp bool, timeStampFormat string) {
	var writer io.Writer
	if exportPath != "" {
		f, err := os.Create(exportPath)
		if err != nil {
			panic(err)
		}
		writer = io.MultiWriter(os.Stdout, f)
	} else {
		writer = os.Stdout
	}
	l._logger = slog.New(slog.NewTextHandler(writer, nil))
}

func (l *ULogger) Info(msg string) {
	if l._logger == nil {
		panic("logger not initialized")
	}
	l._logger.InfoContext(context.Background(), msg)
}

func (l *ULogger) Debug(msg string) {
	if l._logger == nil {
		panic("logger not initialized")
	}
	l._logger.DebugContext(context.Background(), msg)
}

func (l *ULogger) Warn(msg string) {
	if l._logger == nil {
		panic("logger not initialized")
	}
	l._logger.WarnContext(context.Background(), msg)
}

func (l *ULogger) Err(msg string) {
	if l._logger == nil {
		panic("logger not initialized")
	}
	slog.ErrorContext(context.Background(), msg)
}
