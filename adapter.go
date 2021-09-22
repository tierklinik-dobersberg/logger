package logger

import (
	"fmt"
	"log"
	"sync"
	"time"
)

// Severity defines the severity of a log message.
type Severity int

const (
	// Info is used for almost all non-critical log messages.
	Info = Severity(5)
	// Error is used to log error conditions that have not been handled
	// or failed to be handled correctly.
	Error = Severity(0)
)

func (severity Severity) String() string {
	switch severity {
	case Info:
		return "info(5)"
	case Error:
		return "error(0)"
	default:
		return fmt.Sprintf("(%d)", severity)
	}
}

type (
	// Adapter is used to actually write log messages to some final destination.
	Adapter interface {
		// Write is called for each log message and should persist the log message and
		// it's fields somewhere. Write is not allowed to manipulated the Fields map as it
		// may be used for other messages concurrently.
		Write(clock time.Time, severity Severity, msg string, fields Fields)
	}

	// AdapterFunc is convenience type for creating log adapters.
	AdapterFunc func(time.Time, Severity, string, Fields)
)

func (fn AdapterFunc) Write(clock time.Time, severity Severity, msg string, fields Fields) {
	fn(clock, severity, msg, fields)
}

// StdlibAdapter is a log adapter that uses the log package from the standart library.
type StdlibAdapter struct{}

func (*StdlibAdapter) Write(_ time.Time, severity Severity, msg string, fields Fields) {
	for k, v := range fields {
		msg += fmt.Sprintf(" %s=%v", k, v)
	}

	log.Println(msg)
}

var (
	defaultAdapter     Adapter
	defaultAdapterOnce sync.Once
)

// SetDefaultAdapter sets the default logging adapter used by the default logger.
// Note that SetDefaultAdapter can only be called once and must be called before
// any call to DefaultAdapter() or DefaultLogger().
func SetDefaultAdapter(a Adapter) {
	defaultAdapterOnce.Do(func() {
		defaultAdapter = a
	})
}

// DefaultAdapter returns the default logging adapter. If SetDefaultAdapter has not
// been called the default adapter is set to a new StdlibAdapter. Further calls to
// SetDefaultAdapter will be no-ops.
func DefaultAdapter() Adapter {
	defaultAdapterOnce.Do(func() {
		defaultAdapter = new(StdlibAdapter)
	})
	return defaultAdapter
}
