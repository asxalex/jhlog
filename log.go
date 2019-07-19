package log

import (
	"log"
	"os"
	"strings"

	"github.com/op/go-logging"
)

const (
	defaultLevel = logging.DEBUG
)

var (
	Separator = string(os.PathSeparator)
)

// Files is a reference to the opened files

var logMap = make(map[string]*logging.Logger)
var logLevel logging.Level
var format = logging.MustStringFormatter(
	`%{color}%{time:2006-01-02 15:04:05.000} %{shortfunc} > %{level:.4s} %{color:reset} %{message}`,
)
var defaultLogPath = "."

func getFile(logname string) (*File, error) {
	logpath := defaultLogPath
	filename := strings.TrimRight(logpath, Separator) + Separator + logname
	f, err := NewFile(filename, "2006-01-02")
	if err != nil {
		return nil, err
	}
	f.SetRotate(1)
	f.SetAutoDelete(7)
	return f, nil
}

// GetLog returns the logging.Logger through the name
// if the name is the same, it will get the same logger
func GetLog(logname string) *logging.Logger {
	if logger, ok := logMap[logname]; ok {
		return logger
	}
	logger := logging.MustGetLogger(logname)
	var loglevel = defaultLevel

	// f, err := os.OpenFile(strings.TrimRight(logpath, "/")+"/"+logname+".log", os.O_RDWR|os.O_APPEND|os.O_CREATE, 0666)
	f, err := getFile(logname)
	if err != nil {
		log.Fatal(err)
	}
	backend := logging.NewLogBackend(f, "", 0)
	leveledbackendFormatted := logging.NewBackendFormatter(backend, format)
	leveledbackend := logging.AddModuleLevel(leveledbackendFormatted)
	leveledbackend.SetLevel(loglevel, "")
	logger.SetBackend(leveledbackend)
	logMap[logname] = logger
	return logger
}
