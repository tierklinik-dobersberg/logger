package log

import (
	"context"
	"fmt"
	"sync"
	"time"
)

type Fields map[string]interface{}

type Logger interface {
	Info(msg string)
	Infof(msg string, args ...interface{})

	Error(msg string)
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

func From(ctx context.Context) Logger {
	if l, ok := ctx.Value(loggerKey).(Logger); ok {
		return l
	}

	return DefaultLogger()
}

func With(ctx context.Context, log Logger) context.Context {
	return context.WithValue(ctx, loggerKey, log)
}

func ContextFields(ctx context.Context) Fields {
	if f, ok := ctx.Value(fieldsKey).(Fields); ok {
		return f
	}

	return nil
}

func WithFields(ctx context.Context, fields Fields) context.Context {
	newFields := mergeFields(ContextFields(ctx), fields)
	return context.WithValue(ctx, fieldsKey, newFields)
}

func DefaultLogger() Logger {
	defaultLoggerOnce.Do(func() {
		defaultLogger = &logger{
			adapter: DefaultAdapter(),
		}
	})
	return defaultLogger
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
