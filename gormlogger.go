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

type ElapsedTimeFunc func(time.Duration)
type SqlRowFunc func(string, int64)
type ErrorFunc func(error)

type Logging interface {
	Tracef(msg string, data ...interface{})
	Infof(msg string, data ...interface{})
	Warnf(msg string, data ...interface{})
	Errorf(msg string, data ...interface{})
}

func NewLogger(log Logging) *Logger {
	return &Logger{
		Logging: log,
	}
}

type Logger struct {
	Logging

	elapsedTimeFunc ElapsedTimeFunc
	sqlRowFunc      SqlRowFunc
	errorFunc       ErrorFunc
}

func (l *Logger) SetElapsedTimeFunc(fn ElapsedTimeFunc) *Logger {
	l.elapsedTimeFunc = fn
	return l
}

func (l *Logger) SetSqlRowFunc(fn SqlRowFunc) *Logger {
	l.sqlRowFunc = fn
	return l
}

func (l *Logger) SetErrorFunc(fn ErrorFunc) *Logger {
	l.errorFunc = fn
	return l
}

func (l *Logger) LogMode(level logger.LogLevel) logger.Interface {
	return l
}

func (l *Logger) Trace(ctx context.Context, begin time.Time, fc func() (string, int64), err error) {
	elapsed := time.Since(begin)
	if l.elapsedTimeFunc != nil {
		l.elapsedTimeFunc(elapsed)
	}

	sql, rows := fc()
	if rows == -1 {
		if err != nil {
			l.Logging.Errorf(errStr, err, float64(elapsed.Nanoseconds())/1e6, "-", sql)
		} else {
			l.Logging.Tracef(traceStr, float64(elapsed.Nanoseconds())/1e6, "-", sql)
		}
	} else {
		if err != nil {
			l.Logging.Errorf(errStr, err, float64(elapsed.Nanoseconds())/1e6, rows, sql)
		} else {
			l.Logging.Tracef(traceStr, float64(elapsed.Nanoseconds())/1e6, rows, sql)
		}
	}
	if l.sqlRowFunc != nil {
		l.sqlRowFunc(sql, rows)
	}
	if l.errorFunc != nil && err != nil {
		l.errorFunc(err)
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
