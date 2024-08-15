package models

import "sync"

type ReceiverMessageQueue interface {
	ReceiveMessages(wg *sync.WaitGroup) error
}

type LogLevel string

const (
	ErrorLevel LogLevel = "error"
	WarnLevel  LogLevel = "warn"
	DebugLevel LogLevel = "debug"
	InfoLevel  LogLevel = "info"
)

type Logger interface {
	Log(logLevel LogLevel, message string, log ...interface{})
}
