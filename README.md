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
	jhlog.SetDefaultLogPath("./logs")
	jhlog.SetLogLevel(jhlog.INFO)
	logger = jhlog.GetLog("global-log")
}

func main() {
	logger.Debugf("hello world")
	logger.Infof("hello world")
	logger.Warningf("hello world")
	logger.Errorf("hello world")
	logger.Criticalf("hello world")
}
```
