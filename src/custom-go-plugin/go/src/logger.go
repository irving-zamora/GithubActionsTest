package main

import (
	"fmt"
)

type LogLevel int

const (
	Debug LogLevel = iota
	Info
	Error
)

var currentLogLevel = Info // Default log level

func SetLogLevel(level LogLevel) {
	DebugLog("current log level: ", level)
	currentLogLevel = level
}

func log(logLevel LogLevel, format string, args ...interface{}) {
	if logLevel >= currentLogLevel {
		message := format
		if len(args) > 0 {
			message = fmt.Sprintf(format, args...)
		}
		fmt.Println(message)
	}
}

func DebugLog(format string, args ...interface{}) {
	log(Debug, "[DEBUG] "+format, args...)
}

func InfoLog(format string, args ...interface{}) {
	log(Info, "[INFO] "+format, args...)
}

func ErrorLog(format string, args ...interface{}) {
	log(Error, "[ERROR] "+format, args...)
}
