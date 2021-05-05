package main

import (
	"fmt"

	log "github.com/pion/ion-log"
)

var (
	logger = log.NewLogger(log.DebugLevel, "example")
)

func main() {
	logger.Debug("Hello ION!")

	// Change logs level to Info, logger.Debug will nothing output.
	log.SetLogLevel("example", log.InfoLevel)
	logger.Debug("nothing output!")

	logger.Info("Print Info!")

	loggers := log.GetLoggers()

	for _, logger := range loggers {
		fmt.Printf("logger prefix = %v\n", logger.Prefix())
	}
}
