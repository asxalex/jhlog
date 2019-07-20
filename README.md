jhlog
===
log using go-logging

## example

```golang
package main

import (
	"github.com/asxalex/jhlog"
)

var logger *log.Logger

func init() {
	log.SetDefaultLogPath("./logs")
	logger = log.GetLog("global-log")
}

func main() {
	logger.Debugf("hello world")
}
```
