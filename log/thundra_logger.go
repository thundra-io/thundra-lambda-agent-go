package log

import (
	"errors"
	"fmt"
	"log"
	"os"
	"runtime"
	"strings"

	"github.com/thundra-io/thundra-lambda-agent-go/tracer"

	ot "github.com/opentracing/opentracing-go"
	"github.com/thundra-io/thundra-lambda-agent-go/plugin"
)

var (
	logManager *thundraLogManager
	// Logger used as main thundra logger
	Logger     *thundraLogger
	// If we use prebuilt logger functions these are print,panic or fatal
	// we have too add an additional calldepth for our wrapper.
	// It is zero for other functions: trace, debug, info, warn, error.
	additionalCalldepth int
	logLevelCode        int
	spanID              string
)

func init() {
	logManager = &thundraLogManager{}
	Logger = newThundraLogger(logManager)
	logLevelCode = getLogLevelCode()
}

type thundraLogger struct {
	*log.Logger
}

type thundraLogManager struct {
	logs               []*monitoringLog
	recentLogLevel     string // recentLogLevel saves the level of the last log call
	recentLogLevelCode int    // recentLogLevelCode saves the level code of the last log call
}

func newThundraLogger(t *thundraLogManager) *thundraLogger {
	return &thundraLogger{
		Logger: log.New(t, "", 0),
	}
}

// Trace prints trace level logs to logger.
func (l *thundraLogger) Trace(v ...interface{}) {
	if logLevelCode > traceLogLevelCode {
		return
	}
	logManager.recentLogLevel = traceLogLevel
	logManager.recentLogLevelCode = traceLogLevelCode
	l.Output(2, fmt.Sprint(v...))
}

func (l *thundraLogger) TraceWithSpan(span ot.Span, v ...interface{}) {
	setSpanID(span)
	additionalCalldepth = 1
	l.Trace(v...)
}

// Debug prints debug level logs to logger.
func (l *thundraLogger) Debug(v ...interface{}) {
	if logLevelCode > debugLogLevelCode {
		return
	}
	logManager.recentLogLevel = debugLogLevel
	logManager.recentLogLevelCode = debugLogLevelCode
	l.Output(2, fmt.Sprint(v...))
}

func (l *thundraLogger) DebugWithSpan(span ot.Span, v ...interface{}) {
	setSpanID(span)
	additionalCalldepth = 1
	l.Debug(v...)
}

// Info prints info level logs to logger.
func (l *thundraLogger) Info(v ...interface{}) {
	if logLevelCode > infoLogLevelCode {
		return
	}
	logManager.recentLogLevel = infoLogLevel
	logManager.recentLogLevelCode = infoLogLevelCode
	l.Output(2, fmt.Sprint(v...))
}

func (l *thundraLogger) InfoWithSpan(span ot.Span, v ...interface{}) {
	setSpanID(span)
	additionalCalldepth = 1
	l.Info(v...)
}

// Warn prints warn level logs to logger.
func (l *thundraLogger) Warn(v ...interface{}) {
	if logLevelCode > warnLogLevelCode {
		return
	}
	logManager.recentLogLevel = warnLogLevel
	logManager.recentLogLevelCode = warnLogLevelCode
	l.Output(2, fmt.Sprint(v...))
}

func (l *thundraLogger) WarnWithSpan(span ot.Span, v ...interface{}) {
	setSpanID(span)
	additionalCalldepth = 1
	l.Warn(v...)
}

// Error prints error level logs to logger.
func (l *thundraLogger) Error(v ...interface{}) {
	if logLevelCode > errorLogLevelCode {
		return
	}
	logManager.recentLogLevel = errorLogLevel
	logManager.recentLogLevelCode = errorLogLevelCode
	l.Output(2, fmt.Sprint(v...))
}

func (l *thundraLogger) ErrorWithSpan(span ot.Span, v ...interface{}) {
	setSpanID(span)
	additionalCalldepth = 1
	l.Error(v...)
}

// Below are the wrapper functions for standard library's logger.

// Printf sets log level to info and calls standard library's Printf function.
func (l thundraLogger) Printf(format string, v ...interface{}) {
	if logLevelCode > infoLogLevelCode {
		return
	}
	logManager.recentLogLevel = infoLogLevel
	logManager.recentLogLevelCode = infoLogLevelCode
	additionalCalldepth = 1
	l.Logger.Printf(format, v...)
}

// Print sets log level to info and calls standard library's Print function.
func (l thundraLogger) Print(v ...interface{}) {
	if logLevelCode > infoLogLevelCode {
		return
	}
	logManager.recentLogLevel = infoLogLevel
	logManager.recentLogLevelCode = infoLogLevelCode
	additionalCalldepth = 1
	l.Logger.Print(v...)
}

// Println sets log level to info and calls standard library's Println function.
func (l thundraLogger) Println(v ...interface{}) {
	if logLevelCode > infoLogLevelCode {
		return
	}
	logManager.recentLogLevel = infoLogLevel
	logManager.recentLogLevelCode = infoLogLevelCode
	additionalCalldepth = 1
	l.Logger.Println(v...)
}

// Panicf sets log level to error and calls standard library's Panicf function.
func (l thundraLogger) Panicf(format string, v ...interface{}) {
	if logLevelCode > errorLogLevelCode {
		return
	}
	logManager.recentLogLevel = errorLogLevel
	logManager.recentLogLevelCode = errorLogLevelCode
	additionalCalldepth = 1
	l.Logger.Panicf(format, v...)
}

// Panic sets log level to error and calls standard library's Panic function.
func (l thundraLogger) Panic(v ...interface{}) {
	if logLevelCode > errorLogLevelCode {
		return
	}
	logManager.recentLogLevel = errorLogLevel
	logManager.recentLogLevelCode = errorLogLevelCode
	additionalCalldepth = 1
	l.Logger.Panic(v...)
}

// Panicln sets log level to error and calls standard library's Panicln function.
func (l thundraLogger) Panicln(v ...interface{}) {
	if logLevelCode > errorLogLevelCode {
		return
	}
	logManager.recentLogLevel = errorLogLevel
	logManager.recentLogLevelCode = errorLogLevelCode
	additionalCalldepth = 1
	l.Logger.Panicln(v...)
}

// Write stores the log into logs array which will later be used to send monitoredLogs to Thundra collector.
func (t *thundraLogManager) Write(p []byte) (n int, err error) {
	// We need to skip last 3 frames and additionalCalldepth for wrapper functions
	pc, _, line, ok := runtime.Caller(3 + additionalCalldepth)
	contextName := runtime.FuncForPC(pc).Name()

	// Reset addtional calldepth
	additionalCalldepth = 0

	if !ok {
		contextName = "???"
		line = 0
	}

	mL := &monitoringLog{
		logMessage:     string(p),
		logContextName: fmt.Sprintf("%s: %d", contextName, line),
		logTimestamp:   plugin.GetTimestamp(),
		logLevel:       t.recentLogLevel,
		logLevelCode:   t.recentLogLevelCode,
		spanID:         spanID,
	}

	// Reset span id
	spanID = ""

	// Append log
	t.logs = append(t.logs, mL)
	return len(p), nil
}

func (t *thundraLogManager) clearLogs() {
	t.logs = nil
}

func getLogLevelCode() int {
	l := os.Getenv(thundraLogLogLevel)
	thundraLogLevel := strings.ToUpper(l)
	if thundraLogLevel == traceLogLevel {
		return traceLogLevelCode
	} else if thundraLogLevel == debugLogLevel {
		return debugLogLevelCode
	} else if thundraLogLevel == infoLogLevel {
		return infoLogLevelCode
	} else if thundraLogLevel == warnLogLevel {
		return warnLogLevelCode
	} else if thundraLogLevel == errorLogLevel {
		return errorLogLevelCode
	} else if thundraLogLevel == noneLogLevel {
		// Logging is disabled. None of the logs will be sent.
		return noneLogLevelCode
	} else if thundraLogLevel == "" {
		// logLevel is not set, thundra will report all logs.
		return 0
	}

	log.Print(errors.New("invalid " + thundraLogLogLevel + ". Logs are disabled. Use trace, debug, info, warn, error or none."))
	return noneLogLevelCode
}

func setSpanID(span ot.Span) {
	spanCtx, ok := span.Context().(tracer.SpanContext)
	if ok {
		spanID = spanCtx.SpanID
	}
}
