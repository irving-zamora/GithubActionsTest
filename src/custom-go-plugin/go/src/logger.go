// Log function that allows to set what gets logged or not based on the api definition config,
// having Info, Debug and Error levels
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

// SetLogLevel sets the current log level.
// levels are:
// Debug = 0
// Info = 1
// Error = 2
func SetLogLevel(level LogLevel) {
	DebugLog("current log level: ", level)
	currentLogLevel = level
}

// is a helper function that handles actual logging.
// function checks if the given logLevel is greater than or equal to the currentLogLevel,
// and if it is, the message is printed.
func log(logLevel LogLevel, format string, args ...interface{}) {
	if logLevel >= currentLogLevel {
		message := format
		if len(args) > 0 {
			message = fmt.Sprintf(format, args...)
		}
		fmt.Println(message)
	}
}

// DebugLog logs a debug-level message.
func DebugLog(format string, args ...interface{}) {
	log(Debug, "[DEBUG] "+format, args...)
}

// InfoLog logs an info-level message.
func InfoLog(format string, args ...interface{}) {
	log(Info, "[INFO] "+format, args...)
}

// ErrorLog logs an error-level message.
func ErrorLog(format string, args ...interface{}) {
	log(Error, "[ERROR] "+format, args...)
}
