package log

import (
	"testing"
)

var (
	logger = NewLogger(InfoLevel, "test")
)

func TestLogFormat(t *testing.T) {
	logger.Info("xxx")
}
