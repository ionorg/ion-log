package log

import "testing"

func TestLogFormat(t *testing.T) {
	fixByFile := []string{"asm_amd64.s", "proc.go"}
	fixByFunc := []string{}
	Init("debug", fixByFile, fixByFunc)
	Infof("Hello %s!", "ION")
}

func TestLogFixByFunc(t *testing.T) {
	fixByFile := []string{"asm_amd64.s", "proc.go"}
	fixByFunc := []string{"tRunner"}
	Init("debug", fixByFile, fixByFunc)
	Infof("Hello %s!", "ION")
}

func TestLogFixByFile(t *testing.T) {
	fixByFile := []string{"asm_amd64.s", "proc.go"}
	fixByFunc := []string{}
	Init("debug", fixByFile, fixByFunc)
	printOK := make(chan struct{})
	go func() {
		Infof("Hello %s!", "ION")
		printOK <- struct{}{}
	}()
	<-printOK
}
