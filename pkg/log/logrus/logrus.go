package logrus

import (
	"fmt"
	"os"
	"runtime"
	lgrus "github.com/sirupsen/logrus"
	"github.com/maxott/magda-cli/pkg/log"
)

// Logger struct wrapping around a logrus logger.
type logger struct {
	lr *lgrus.Logger
	tags lgrus.Fields
}

// Debug is used for logging debug level logs.
func (l *logger) Debug(msg string) {
	l.log(lgrus.DebugLevel, msg)
}

// Debugf is used for logging debug level logs.
func (l *logger) Debugf(fmt string, args ...interface{}) {
	l.log(lgrus.DebugLevel, fmt, args...)
}

// Info is used for logging info level logs.
func (l *logger) Info(msg string) {
	l.log(lgrus.InfoLevel, msg)
}

// Infof is used for logging info level logs.
func (l *logger) Infof(fmt string, args ...interface{}) {
	l.log(lgrus.InfoLevel, fmt, args...)
}

// Warn is used for logging warn level logs.
func (l *logger) Warn(msg string) {
	l.log(lgrus.WarnLevel, msg)
}

// Warnf is used for logging warning level logs.
func (l *logger) Warnf(fmt string, args ...interface{}) {
	l.log(lgrus.WarnLevel, fmt, args...)
}

// Error is used for logging error level logs.
func (l *logger) Error(err error, msg string) error {
	l.log(lgrus.ErrorLevel, msg)
	if err == nil {
		err = fmt.Errorf(msg)
	}
	return err
}

// Errorf is used for logging error level logs.
func (l *logger) Errorf(err error, format string, args ...interface{}) error {
	l.log(lgrus.ErrorLevel, format, args...)
	if err == nil {
		err = fmt.Errorf(format, args...)
	}
	return err
}

// Fatal is used for logging fatal level logs.
func (l *logger) Fatal(msg string) {
	l.log(lgrus.FatalLevel, msg)
}

// Fatals is used for logging fatal level logs.
func (l *logger) Fatalf(fmt string, args ...interface{}) {
	l.log(lgrus.FatalLevel, fmt, args...)
}

func (l *logger) WithTags(tags log.Tags) log.Logger {
	lt := l.tags
	if lt == nil {
		lt = lgrus.Fields{}
	}
	for k, v := range tags {
		lt[k] = v
	}
	return &logger{lr: l.lr, tags: lt}
}

func (l *logger) With(key string, value interface{}) log.Logger {
	tags := l.tags
	if tags == nil {
		tags = lgrus.Fields{}
	}
	tags[key] = value
	return &logger{lr: l.lr, tags: tags}
}

func (l *logger) log(level lgrus.Level, format string, args ...interface{}) {
	if !l.lr.IsLevelEnabled(level) {
		return
	}

	tags := l.tags
	if tags == nil {
		tags = lgrus.Fields{}
	}

	pcs := make([]uintptr, 1)
	// there should always be just one local function above the callstack
	_ = runtime.Callers(3, pcs)
	tags["__frame"] = pcs

	e := l.lr.WithFields(tags)
	e.Logf(level, format, args...)
}

// NewDefaultLogger creates a new instance of logrus logger.
func NewSimpleLogger(level lgrus.Level) log.Logger {
	l := lgrus.Logger{
		Out:          os.Stderr,
		Formatter:    new(lgrus.TextFormatter),
		Hooks:        make(lgrus.LevelHooks),
		Level:        level, // lgrus.InfoLevel,
		ExitFunc:     os.Exit,
		ReportCaller: false,
	}
	return NewLogger(&l)
}

// NewLogger creates a new instance of logrus logger.
func NewLogger(lrl *lgrus.Logger) log.Logger {
	hook := new(CallerHook)
	lrl.Hooks.Add(hook)
	return &logger{lr: lrl}
}

// Add fixed frame
type CallerHook struct {}

func (hook *CallerHook) Fire(entry *lgrus.Entry) error {
	f := entry.Data
	pcs := f["__frame"].([]uintptr)
	delete(f, "__frame")
	frames := runtime.CallersFrames(pcs)
	frame, _ := frames.Next()
	entry.Caller = &frame
	return nil
}

func (hook *CallerHook) Levels() []lgrus.Level {
	return []lgrus.Level{
		lgrus.TraceLevel,
		lgrus.DebugLevel,
		lgrus.InfoLevel,
		lgrus.WarnLevel,
		lgrus.ErrorLevel,
		lgrus.FatalLevel,
		lgrus.PanicLevel,
	}
}
