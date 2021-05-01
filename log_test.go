package log

import "testing"

func TestLogFormat(t *testing.T) {
	Init("debug")
	Infof("Hello %s!", "ION")
}
