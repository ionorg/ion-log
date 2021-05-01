package main

import (
	log "github.com/pion/ion-log"
)

func init() {
	log.Init("debug")
}

func main() {
	log.Infof("Hello ION!")
}
