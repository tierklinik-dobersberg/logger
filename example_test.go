package logger_test

import "github.com/tierklinik-dobersberg/logger"

func Example() {
	logger.SetDefaultAdapter(new(logger.StdlibAdapter))

	logger.DefaultLogger().Infof("Some log message")

	fieldLogger := logger.DefaultLogger().WithFields(logger.Fields{
		"name": "main",
		"user": "demo",
	})

	fieldLogger.Error("another log message")
	fieldLogger.WithFields(logger.Fields{"session": "test"}).Infof("once more")
}
