package logger

import (
	"context"
	"os"
)

// Infof logs an info message with the logger and fields from
// ctx (or the DefaultLogger).
func Infof(ctx context.Context, msg string, args ...interface{}) {
	From(ctx).WithFields(ContextFields(ctx)).Infof(msg, args...)
}

// Errorf logs an error message with the logger and fields from
// ctx (or the DefaultLogger).
func Errorf(ctx context.Context, msg string, args ...interface{}) {
	From(ctx).WithFields(ContextFields(ctx)).Errorf(msg, args...)
}

// Fatalf is like Errorf but exits the process.
func Fatalf(ctx context.Context, msg string, args ...interface{}) {
	Errorf(ctx, msg, args...)
	os.Exit(1)
}
