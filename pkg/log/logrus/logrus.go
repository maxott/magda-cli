package logrus

import (
	"fmt"
	"os"
	lgrus "github.com/sirupsen/logrus"
	"github.com/maxott/magda-cli/pkg/log"
)

// Logger struct wrapping around a logrus logger.
type logger struct {
	lr *lgrus.Logger
	tags *lgrus.Fields
}

// Debug is used for logging debug level logs.
func (l *logger) Debug(msg string) {
	l.lr.Debug(msg)
}

// Debugf is used for logging debug level logs.
func (l *logger) Debugf(fmt string, args ...interface{}) {
	l.lr.Debugf(fmt, args...)
}

// Info is used for logging info level logs.
func (l *logger) Info(msg string) {
	l.lr.Info(msg)
}

// Infof is used for logging info level logs.
func (l *logger) Infof(fmt string, args ...interface{}) {
	l.lr.Infof(fmt, args...)
}

// Warn is used for logging warn level logs.
func (l *logger) Warn(msg string) {
	l.lr.Warn(msg)
}

// Warnf is used for logging warning level logs.
func (l *logger) Warnf(fmt string, args ...interface{}) {
	l.lr.Warnf(fmt, args...)
}

// Error is used for logging error level logs.
func (l *logger) Error(err error, msg string) error {
	l.lr.Error(msg)
	if err == nil {
		err = fmt.Errorf(msg)
	}
	return err
}

// Errorf is used for logging error level logs.
func (l *logger) Errorf(err error, format string, args ...interface{}) error {
	l.lr.Errorf(format, args...)
	if err == nil {
		err = fmt.Errorf(format, args...)
	}
	return err
}

// Fatal is used for logging fatal level logs.
func (l *logger) Fatal(msg string) {
	l.lr.Fatal(msg)
}

// Fatals is used for logging fatal level logs.
func (l *logger) Fatalf(fmt string, args ...interface{}) {
	l.lr.Fatalf(fmt, args...)
}

func (l *logger) WithTags(tags log.Tags) log.Logger {
	lt := l.tags
	if lt == nil {
		lt = &lgrus.Fields{}
	}
	for k, v := range tags {
		(*lt)[k] = v
	}
	return &logger{lr: l.lr, tags: lt}
}

func (l *logger) With(key string, value interface{}) log.Logger {
	tags := l.tags
	if tags == nil {
		tags = &lgrus.Fields{}
	}
	(*tags)[key] = value
	return &logger{lr: l.lr, tags: tags}
}

// NewDefaultLogger creates a new instance of logrus logger.
func NewDefaultLogger() log.Logger {
	l := &lgrus.Logger{
		Out:          os.Stderr,
		Formatter:    new(lgrus.TextFormatter),
		Hooks:        make(lgrus.LevelHooks),
		Level:        lgrus.DebugLevel, // lgrus.InfoLevel,
		ExitFunc:     os.Exit,
		ReportCaller: false,
	}
	return &logger{lr: l}
}

// NewLogger creates a new instance of logrus logger.
func NewLogger(lrl *lgrus.Logger) log.Logger {
	return &logger{lr: lrl}
}
