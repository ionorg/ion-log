package main

import (
	"fmt"

	log "github.com/pion/ion-log"
)

var (
	logger = log.NewLoggerWithFields(log.DebugLevel, "example", log.Fields{"field1": "value1", "field2": "value2"})
)

func main() {
	logger.Debug("Hello ION!")

	// Change logs level to Info, logger.Debug will nothing output.
	log.SetLogLevel("example", log.InfoLevel)
	logger.Debug("nothing output!")

	logger.Info("Info!")

	logger.Warn("Warn!")

	logger.Error("Error!")

	loggers := log.GetLoggers()

	for _, logger := range loggers {
		fmt.Printf("logger prefix = %v\n", logger.Prefix())
	}
}
