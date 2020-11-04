package main

import (
	log "github.com/pion/ion-log"
)

func init() {
	fixByFile := []string{"asm_amd64.s", "proc.go"}
	fixByFunc := []string{}
	log.Init("debug", fixByFile, fixByFunc)
}

func main() {
	log.Infof("Hello ION!")
}
