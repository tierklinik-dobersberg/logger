package logger

import "time"

type multiAdapter struct {
	adapters []Adapter
}

func (multi *multiAdapter) Write(clock time.Time, severity Severity, msg string, fields Fields) {
	for _, a := range multi.adapters {
		a.Write(clock, severity, msg, fields)
	}
}

func MultiAdapter(adapters ...Adapter) Adapter {
	return &multiAdapter{adapters}
}
