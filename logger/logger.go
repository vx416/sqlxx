package logger

import (
	"context"
	"fmt"
	"time"
)

type (
	Logger interface {
		Warn(s string, fields map[string]interface{})
		Info(s string, fields map[string]interface{})
		Debug(s string, fields map[string]interface{})
		Error(s string, fields map[string]interface{})
	}

	LoggerKey struct{}
	LevelKey  struct{}
	Level     uint8
)

const (
	Debug Level = iota + 1
	Info
	Off
)

var (
	LogLevel      = Debug
	SlowThreshold = 1 * time.Second
	Colorful      = false
)

func AttachLogger(ctx context.Context, log Logger) context.Context {
	return context.WithValue(ctx, LoggerKey{}, log)
}

func WithLogLevel(ctx context.Context, level Level) context.Context {
	return context.WithValue(ctx, LevelKey{}, level)
}

func getLevel(ctx context.Context) Level {
	level, ok := ctx.Value(LevelKey{}).(Level)
	if !ok {
		return LogLevel
	}
	return level
}

func GetLogger(ctx context.Context) Logger {
	logger, ok := ctx.Value(LoggerKey{}).(Logger)
	if !ok {
		return nil
	}
	return logger
}

const (
	colorBlack = iota + 30
	colorRed
	colorGreen
	colorYellow
	colorBlue
	colorMagenta
	colorCyan
	colorWhite

	colorBold     = 1
	colorDarkGray = 90
)

func Print(ctx context.Context, rows int64, err error, cost time.Duration, query string, args ...interface{}) {
	level := getLevel(ctx)
	if level == Off {
		return
	}

	l := GetLogger(ctx)
	if l == nil {
		return
	}

	sql := ExplainSQL(query, args...)
	if Colorful {
		sql = colorize(sql, colorGreen)
	}
	fields := map[string]interface{}{
		"rows_affected": rows,
		"db_cost":       cost,
	}
	if err != nil {
		fields["error"] = err.Error()
	}

	if err != nil {
		l.Error(sql, fields)
	} else if cost > SlowThreshold {
		l.Warn(sql, fields)
	} else if level == Info {
		l.Info(sql, fields)
	} else {
		l.Debug(sql, fields)
	}
}

func colorize(s string, c int) string {
	return fmt.Sprintf("\x1b[%dm%s\x1b[0m", c, s)
}
