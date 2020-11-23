package log

import (
	"fmt"
	"log"
	"sync"
	"time"
)

type Severity string

const (
	Info  = Severity("info")
	Error = Severity("error")
)

type (
	Adapter interface {
		Write(clock time.Time, severity Severity, msg string, fields Fields)
	}

	AdapterFunc func(Severity, string, Fields)
)

func (fn AdapterFunc) Write(clock time.Time, severity Severity, msg string, fields Fields) {
	fn(clock, severity, msg, fields)
}

type StdlibAdapter struct{}

func (*StdlibAdapter) Write(_ time.Time, severity Severity, msg string, fields Fields) {
	for k, v := range fields {
		msg += fmt.Sprintf(" %s=%q", k, v)
	}

	log.Println(msg)
}

var (
	defaultAdapter     Adapter
	defaultAdapterOnce sync.Once
)

func SetDefaultAdapter(a Adapter) {
	defaultAdapterOnce.Do(func() {
		defaultAdapter = a
	})
}

func DefaultAdapter() Adapter {
	defaultAdapterOnce.Do(func() {
		defaultAdapter = new(StdlibAdapter)
	})
	return defaultAdapter
}
