package logger

import (
	"context"
	"fmt"
	"sync"
	"time"
)

// Fields is a map of additional fields used for logging.
type Fields map[string]interface{}

// Logger describes the general log interfaces exposed by this
// package.
type Logger interface {
	// Info logs a message with info level.
	Info(msg string)
	// Infof logs a format message.
	Infof(msg string, args ...interface{})
	// Error logs a message with error level.
	Error(msg string)
	// Errorf logs a format message with error level.
	Errorf(msg string, args ...interface{})
	// WithFields returns a new logger that also logs fields.
	WithFields(Fields) Logger
}

var (
	fieldsKey = struct{}{}
	loggerKey = struct{}{}

	defaultLogger     Logger
	defaultLoggerOnce sync.Once
)

// From returns the logger associated with ctx or the
// default logger.
func From(ctx context.Context) Logger {
	if l, ok := ctx.Value(loggerKey).(Logger); ok {
		return l
	}

	return DefaultLogger()
}

// With creates a new context from ctx and adds log as the logger.
func With(ctx context.Context, log Logger) context.Context {
	return context.WithValue(ctx, loggerKey, log)
}

// ContextFields returns all log-fields that are associated with
// ctx.
func ContextFields(ctx context.Context) Fields {
	if f, ok := ctx.Value(fieldsKey).(Fields); ok {
		return f
	}

	return nil
}

// WithFields returns a child context that has fields associated and merged
// with any already associated Fields.
func WithFields(ctx context.Context, fields Fields) context.Context {
	newFields := mergeFields(ContextFields(ctx), fields)
	return context.WithValue(ctx, fieldsKey, newFields)
}

// DefaultLogger returns the default logger.
func DefaultLogger() Logger {
	defaultLoggerOnce.Do(func() {
		defaultLogger = &logger{
			adapter: DefaultAdapter(),
		}
	})
	return defaultLogger
}

// New returns a new logger with the given adapter.
func New(adapter Adapter) Logger {
	return &logger{adapter: adapter}
}

type logger struct {
	fields  Fields
	adapter Adapter
}

func (log *logger) Info(msg string) {
	log.adapter.Write(time.Now(), Info, msg, log.fields)
}

func (log *logger) Infof(msg string, args ...interface{}) {
	log.adapter.Write(time.Now(), Info, fmt.Sprintf(msg, args...), log.fields)
}

func (log *logger) Error(msg string) {
	log.adapter.Write(time.Now(), Error, msg, log.fields)
}

func (log *logger) Errorf(msg string, args ...interface{}) {
	log.adapter.Write(time.Now(), Error, fmt.Sprintf(msg, args...), log.fields)
}

func (log *logger) WithFields(fields Fields) Logger {
	newFields := mergeFields(log.fields, fields)
	return &logger{
		fields:  newFields,
		adapter: log.adapter,
	}
}

func mergeFields(a, b Fields) Fields {
	copy := make(Fields, len(a)+len(b))

	for k, v := range a {
		copy[k] = v
	}
	for k, v := range b {
		copy[k] = v
	}

	if len(copy) == 0 {
		return nil
	}

	return copy
}
