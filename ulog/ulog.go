package ulog

import (
	"context"
	"encoding/json"
	"io"
	"log/slog"
	"os"
)

const LOG_LEVEL_INFO = 1
const LOG_LEVEL_DEBUG = 2
const LOG_LEVEL_WARN = 3
const LOG_LEVEL_ERROR = 4

type ULogger struct {
	Level   int
	_tag    string
	_logger *slog.Logger
}

var _ulogger *ULogger = nil

func GetULogger() *ULogger {
	if _ulogger == nil {
		_ulogger = &ULogger{}
	}
	return _ulogger
}

func Log() *ULogger {
	return GetULogger()
}
func (l *ULogger) Tag(tag string) *ULogger {
	l._tag = tag
	return l
}

func Config(level int, exportPath string, trace bool) {
	logger := GetULogger()
	logger.Level = level

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
	var Lv slog.Leveler
	if level == LOG_LEVEL_INFO {
		Lv = slog.LevelInfo
	} else if level == LOG_LEVEL_DEBUG {
		Lv = slog.LevelDebug
	} else if level == LOG_LEVEL_WARN {
		Lv = slog.LevelWarn
	} else if level == LOG_LEVEL_ERROR {
		Lv = slog.LevelError
	} else {
		Lv = slog.LevelDebug
	}

	opts := &slog.HandlerOptions{
		Level:     Lv,
		AddSource: trace,
	}
	logger._logger = slog.New(slog.NewJSONHandler(writer, opts))
}

func (l *ULogger) I(tag string, msg any) {
	if l._logger == nil {
		panic("logger not initialized")
	}

	if _, ok := msg.(string); ok {
		l._logger.Info(deco_tag(tag, msg.(string)))
	} else {
		// log, err := json.MarshalIndent(msg, "", " ")
		// if err != nil {
		// 	l._logger.Error("json marshal failed")
		// }
		log, _ := json.Marshal(msg)

		l._logger.Info(deco_tag(tag, string(log)))
	}

}

/**
 * deprecated
 */
func (l *ULogger) D(tag string, msg any) {
	if l._logger == nil {
		panic("logger not initialized")
	}
	log, _ := json.Marshal(msg)
	l._logger.DebugContext(context.Background(), deco_tag(tag, string(log)))
}

func (l *ULogger) W(tag string, msg any) {
	if l._logger == nil {
		panic("logger not initialized")
	}
	log, err := json.MarshalIndent(msg, "", " ")
	if err != nil {
		l._logger.Error("json marshal failed")
	}
	l._logger.WarnContext(context.Background(), string(log))
	// log, _ := json.Marshal(msg)
	// l._logger.WarnContext(context.Background(), deco_tag(tag, string(log)))
}

func (l *ULogger) E(tag string, msg any) {
	if l._logger == nil {
		panic("logger not initialized")
	}
	//determine if msg is a string or a struct
	if _, ok := msg.(string); ok {
		l._logger.ErrorContext(context.Background(), deco_tag(tag, msg.(string)))
		return
	} else {
		log, _ := json.Marshal(msg)
		l._logger.ErrorContext(context.Background(), deco_tag(tag, string(log)))
	}

}

/**
 * decorate tag
 */
func deco_tag(tag string, msg string) string {
	return "[" + tag + "] " + msg
}
