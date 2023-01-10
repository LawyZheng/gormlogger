package gormlogger

import (
	"context"
	"time"

	"gorm.io/gorm/logger"
)

const (
	traceStr string = "[%.3fms] [rows:%v] %s"
	errStr   string = "[error = %s] [%.3fms] [rows:%v] %s"
)

type FmtLog func(msg string, data ...interface{})

type TraceTripper func(elapsed time.Duration, fc func() (sql string, row int64), e error) (bool, FmtLog)

func newTraceTripper(logger Logging) TraceTripper {
	return func(elapsed time.Duration, fc func() (sql string, row int64), e error) (bool, FmtLog) {
		if e != nil {
			return true, logger.Errorf
		} else {
			return true, logger.Tracef
		}
	}
}

type Logging interface {
	Tracef(msg string, data ...interface{})
	Infof(msg string, data ...interface{})
	Warnf(msg string, data ...interface{})
	Errorf(msg string, data ...interface{})
}

func NewLogger(log Logging) *Logger {
	return &Logger{
		Logging:      log,
		traceTripper: newTraceTripper(log),
	}
}

type Logger struct {
	Logging

	traceTripper TraceTripper
}

func (l *Logger) SetTraceTripper(tripper TraceTripper) *Logger {
	l.traceTripper = tripper
	return l
}

func (l *Logger) LogMode(level logger.LogLevel) logger.Interface {
	return l
}

func (l *Logger) Trace(ctx context.Context, begin time.Time, fc func() (string, int64), err error) {
	if l.traceTripper == nil {
		return
	}

	elapsed := time.Since(begin)

	ok, fn := l.traceTripper(elapsed, fc, err)
	if !ok || fn == nil {
		return
	}

	sql, rows := fc()
	if rows == -1 {
		if err != nil {
			fn(errStr, err, float64(elapsed.Nanoseconds())/1e6, "-", sql)
		} else {
			fn(traceStr, float64(elapsed.Nanoseconds())/1e6, "-", sql)
		}
	} else {
		if err != nil {
			fn(errStr, err, float64(elapsed.Nanoseconds())/1e6, rows, sql)
		} else {
			fn(traceStr, float64(elapsed.Nanoseconds())/1e6, rows, sql)
		}
	}
}

func (l *Logger) Info(ctx context.Context, msg string, data ...interface{}) {
	l.Logging.Infof(msg, data...)
}

func (l *Logger) Warn(ctx context.Context, msg string, data ...interface{}) {
	l.Logging.Warnf(msg, data...)
}

func (l *Logger) Error(ctx context.Context, msg string, data ...interface{}) {
	l.Logging.Errorf(msg, data...)
}
