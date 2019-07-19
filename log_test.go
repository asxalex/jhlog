package log

import (
	"testing"
)

func BenchmarkLog(b *testing.B) {
	logger := GetLog("log_test")
	for i := 0; i < b.N; i++ {
		logger.Debugf("hello world this is a test %d", 1)
		logger.Infof("hello world this is a test2 %d", 2)
		logger.Errorf("hello world this is a test4 %d", 4)
	}
}

/*
func TestLog(t *testing.T) {
	l := GetLog("log_test")
	l.Debug("hello")
	l.Info("hello")
	l.Warning("hello")
	l.Error("hello")
	l.Critical("hello")
}
*/
