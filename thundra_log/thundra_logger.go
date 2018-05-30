package thundra_log

import (
	"fmt"
	"github.com/thundra-io/thundra-lambda-agent-go/plugin"
	"log"
	"runtime"
)

var (
	logManager *thundraLogManager
	Logger     *thundraLogger

	// If we use prebuilt logger functions these are print,panic or fatal
	// we have too add an additional calldepth for our wrapper.
	// It is zero for other functions: trace, debug, info, warn, error.
	additionalCalldepth int
)

func init() {
	logManager = &thundraLogManager{}
	Logger = newThundraLogger(logManager)
}

type thundraLogger struct {
	*log.Logger
}

type thundraLogManager struct {
	logs             []*monitoredLog
	recentLogLevel   string // recentLogLevel saves the level of the last log call
	recentLogLevelId int    // recentLogLevelId saves the level id of the last log call
}

func newThundraLogger(t *thundraLogManager) *thundraLogger {
	//flag := log.Ldate | log.Ltime | log.Lmicroseconds
	return &thundraLogger{
		Logger: log.New(t, "", 0),
	}
}

// Trace prints trace level logs to logger
func (l *thundraLogger) Trace(v ...interface{}) {
	logManager.recentLogLevel = traceLogLevel
	logManager.recentLogLevelId = traceLogLevelId
	l.Output(2, fmt.Sprint(v...))
}

// Debug prints debug level logs to logger
func (l *thundraLogger) Debug(v ...interface{}) {
	logManager.recentLogLevel = debugLogLevel
	logManager.recentLogLevelId = debugLogLevelId
	l.Output(2, fmt.Sprint(v...))
}

// Info prints info level logs to logger
func (l *thundraLogger) Info(v ...interface{}) {
	logManager.recentLogLevel = infoLogLevel
	logManager.recentLogLevelId = infoLogLevelId
	l.Output(2, fmt.Sprint(v...))
}

// Warn prints warn level logs to logger
func (l *thundraLogger) Warn(v ...interface{}) {
	logManager.recentLogLevel = warnLogLevel
	logManager.recentLogLevelId = warnLogLevelId
	l.Output(2, fmt.Sprint(v...))
}

// Error prints error level logs to logger
func (l *thundraLogger) Error(v ...interface{}) {
	logManager.recentLogLevel = errorLogLevel
	logManager.recentLogLevelId = errorLogLevelId
	l.Output(2, fmt.Sprint(v...))
}

// Below are the wrapper functions for standart library's logger.

// Printf sets log level to info and calls standart library's Printf function
func (l thundraLogger) Printf(format string, v ...interface{}) {
	logManager.recentLogLevel = infoLogLevel
	logManager.recentLogLevelId = infoLogLevelId
	additionalCalldepth = 1
	l.Logger.Printf(format, v)
}

// Print sets log level to info and calls standart library's Print function
func (l thundraLogger) Print(v ...interface{}) {
	logManager.recentLogLevel = infoLogLevel
	logManager.recentLogLevelId = infoLogLevelId
	additionalCalldepth = 1
	l.Logger.Print(v)
}

// Println sets log level to info and calls standart library's Println function
func (l thundraLogger) Println(v ...interface{}) {
	logManager.recentLogLevel = infoLogLevel
	logManager.recentLogLevelId = infoLogLevelId
	additionalCalldepth = 1
	l.Logger.Println(v)
}

// Panicf sets log level to error and calls standart library's Panicf function
func (l thundraLogger) Panicf(format string, v ...interface{}) {
	logManager.recentLogLevel = errorLogLevel
	logManager.recentLogLevelId = errorLogLevelId
	additionalCalldepth = 1
	l.Logger.Panicf(format, v)
}

// Panic sets log level to error and calls standart library's Panic function
func (l thundraLogger) Panic(v ...interface{}) {
	logManager.recentLogLevel = errorLogLevel
	logManager.recentLogLevelId = errorLogLevelId
	additionalCalldepth = 1
	l.Logger.Panic(v)
}

// Panicln sets log level to error and calls standart library's Panicln function
func (l thundraLogger) Panicln(v ...interface{}) {
	logManager.recentLogLevel = errorLogLevel
	logManager.recentLogLevelId = errorLogLevelId
	additionalCalldepth = 1
	l.Logger.Panicln(v)
}

// Write stores the log into logs array which will later be used to send monitoredLogs to Thundra collector
func (t *thundraLogManager) Write(p []byte) (n int, err error) {
	// We need to skip last 3 frames and additionalCalldepth for wrapper functions
	_, file, line, ok := runtime.Caller(3 + additionalCalldepth)
	additionalCalldepth = 0 //reset it
	if !ok {
		file = "???"
		line = 0
	}

	mL := &monitoredLog{
		log:          string(p),
		logMessage:   string(p),
		loggerName:   fmt.Sprintf("%s %d", file, line),
		logTimestamp: plugin.GetTimestamp(),
		logLevel:     t.recentLogLevel,
		logLevelId:   t.recentLogLevelId,
	}
	t.logs = append(t.logs, mL)
	return len(p), nil
}

func (t *thundraLogManager) clearLogs() {
	t.logs = nil
}
