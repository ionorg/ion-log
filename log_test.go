package log_test

import (
	"testing"

	log "github.com/pion/ion-log"
)

var (
	logger = log.NewLogger(log.InfoLevel, "test")
)

func TestLogFormat(t *testing.T) {
	logger.Info("logger info print")
	logger.Infof("logger with format %v", "ION")
}

func TestDefaultLogger(t *testing.T) {
	log.Infof("log with format %v", "ION")
}
