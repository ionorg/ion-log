# ion-log
## QuickStart
```
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
```
## Feature
GoodFormat: 
```
[date time]  [Level] [Prefix]: [file:line] [Func] YourLog fields=value
```
```
[2021-05-06 09:54:49.658]  INFO default: [log_test.go:19] [ion-log_test.TestDefaultLogger] log with format ION somefield=foo
```

