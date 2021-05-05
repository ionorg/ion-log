package main

import (
	log "github.com/pion/ion-log"
)

var (
	logger = log.NewLogger(log.DebugLevel, "example")
)

func main() {
	logger.Info("Hello ION!")
}
