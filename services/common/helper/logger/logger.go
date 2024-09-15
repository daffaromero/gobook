package logger

import (
	"fmt"
	"io"
	"log"
	"os"
	"strings"
	"sync"

	"github.com/daffaromero/gobook/services/common/utils"
)

// LogLevel represents the severity of a log message
type LogLevel int

const (
	DEBUG LogLevel = iota
	INFO
	WARN
	ERROR
	FATAL
	PANIC
)

var logLevelNames = map[LogLevel]string{
	DEBUG: "DEBUG",
	INFO:  "INFO",
	WARN:  "WARN",
	ERROR: "ERROR",
	FATAL: "FATAL",
	PANIC: "PANIC",
}

var (
	globalLogLevel LogLevel
	once           sync.Once
)

// Log represents a logger instance
type Log struct {
	prefix string
	stdout *log.Logger
	stderr *log.Logger
}

// New creates a new Log instance
func New(prefix string) *Log {
	once.Do(initGlobalLogLevel)

	return &Log{
		prefix: prefix,
		stdout: log.New(os.Stdout, "", log.Ldate|log.Ltime),
		stderr: log.New(os.Stderr, "", log.Ldate|log.Ltime),
	}
}

func initGlobalLogLevel() {
	levelStr := strings.ToUpper(utils.GetEnv("LOG_LEVEL"))
	for level, name := range logLevelNames {
		if name == levelStr {
			globalLogLevel = level
			return
		}
	}
	globalLogLevel = INFO
}

func (l *Log) log(level LogLevel, w io.Writer, message string, args ...interface{}) {
	if level < globalLogLevel {
		return
	}

	prefix := fmt.Sprintf("[%s][%s] ", logLevelNames[level], l.prefix)
	logger := log.New(w, prefix, log.Ldate|log.Ltime)
	logger.Printf(message, args...)
}

func (l *Log) Debug(message string, args ...interface{}) {
	l.log(DEBUG, os.Stdout, message, args...)
}

func (l *Log) Info(message string, args ...interface{}) {
	l.log(INFO, os.Stdout, message, args...)
}

func (l *Log) Warn(message string, args ...interface{}) {
	l.log(WARN, os.Stdout, message, args...)
}

func (l *Log) Error(message string, args ...interface{}) {
	l.log(ERROR, os.Stderr, message, args...)
}

func (l *Log) Fatal(message string, args ...interface{}) {
	l.log(FATAL, os.Stderr, message, args...)
	os.Exit(1)
}

func (l *Log) Panic(message string, args ...interface{}) {
	l.log(PANIC, os.Stderr, message, args...)
	panic(fmt.Sprintf(message, args...))
}
