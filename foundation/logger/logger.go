// Package logger provides a convenience function to constructing a logger
// for use. This is required not just for applications but for testing.
package logger

import "github.com/sirupsen/logrus"

// New constructs a Logger that writes to stdout and
// provides human-readable timestamps.
func New() *logrus.Logger {

	logger := logrus.New()
	logger.SetLevel(logrus.InfoLevel)
	logger.SetFormatter(&logrus.TextFormatter{
		FullTimestamp:          true,
		TimestampFormat:        "2006-01-02 15:04:05",
		ForceColors:            true,
		DisableLevelTruncation: true,
	})

	return logger
}
