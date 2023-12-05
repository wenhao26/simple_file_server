package config

import (
	"time"
)

var (
	LogFile     string
	LogConsole  bool
	LogLevel    string
	ServerPort  string
	FileStorage string
)

func init() {
	LogFile = "runtime/logs/server.log"
	LogConsole = true
	LogLevel = "debug"
	ServerPort = "8089"
	FileStorage = fileStorage()
}

func IsDev() bool {
	return LogLevel == "debug"
}

func fileStorage() string {
	year := time.Now().Format("2006")
	month := time.Now().Format("01")
	return "files/" + year + "/" + month
}
