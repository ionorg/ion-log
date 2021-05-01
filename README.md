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
[date time][file:line][Level][Func] => YourLog
```
```
[2021-05-01 17:33:47][main.go:12][INFO][main.main] => Hello ION!
```

