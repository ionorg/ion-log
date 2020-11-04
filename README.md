# ion-log
## QuickStart
```
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
```
## Feature
* GoodFormat: [date time] [Level] [Line][File][Func] => YourLog
```
 [2020-11-04 16:13:54.593] [INFO] [14][main.go][main] => Hello ION!
```
* FixByHand: you can fixByFile or fixByFunc when you found the log line not right

