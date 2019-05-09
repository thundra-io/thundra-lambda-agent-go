package log

import (
	"encoding/json"
	"log"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/thundra-io/thundra-lambda-agent-go/utils"
)

const (
	testMessage          = "testMessage"
	expectedTestMessage  = "testMessage\n"
	formattedTestMessage = "[testMessage]\n"
)

func TestLoggerTrace(t *testing.T) {
	Logger.Trace(testMessage)
	assert.Equal(t, traceLogLevel, logManager.recentLogLevel)
	assert.Equal(t, traceLogLevelCode, logManager.recentLogLevelCode)
	assert.Equal(t, expectedTestMessage, logManager.logs[0].logMessage)
	logManager.clearLogs()
}

func TestLoggerDebug(t *testing.T) {
	Logger.Debug(testMessage)
	assert.Equal(t, debugLogLevel, logManager.recentLogLevel)
	assert.Equal(t, debugLogLevelCode, logManager.recentLogLevelCode)
	assert.Equal(t, expectedTestMessage, logManager.logs[0].logMessage)
	logManager.clearLogs()
}

func TestLoggerInfo(t *testing.T) {
	Logger.Info(testMessage)
	assert.Equal(t, infoLogLevel, logManager.recentLogLevel)
	assert.Equal(t, infoLogLevelCode, logManager.recentLogLevelCode)
	assert.Equal(t, expectedTestMessage, logManager.logs[0].logMessage)
	logManager.clearLogs()
}

func TestLoggerWarn(t *testing.T) {
	Logger.Warn(testMessage)
	assert.Equal(t, warnLogLevel, logManager.recentLogLevel)
	assert.Equal(t, warnLogLevelCode, logManager.recentLogLevelCode)
	assert.Equal(t, testMessage+"\n", logManager.logs[0].logMessage)
	logManager.clearLogs()
}

func TestLoggerError(t *testing.T) {
	Logger.Error(testMessage)
	assert.Equal(t, errorLogLevel, logManager.recentLogLevel)
	assert.Equal(t, errorLogLevelCode, logManager.recentLogLevelCode)
	assert.Equal(t, expectedTestMessage, logManager.logs[0].logMessage)
	logManager.clearLogs()
}

func TestLoggerPrintf(t *testing.T) {
	Logger.Printf("[%s]", testMessage)
	assert.Equal(t, infoLogLevel, logManager.recentLogLevel)
	assert.Equal(t, infoLogLevelCode, logManager.recentLogLevelCode)
	assert.Equal(t, formattedTestMessage, logManager.logs[0].logMessage)
	logManager.clearLogs()
}

func TestLoggerPrint(t *testing.T) {
	Logger.Print(testMessage)
	assert.Equal(t, infoLogLevel, logManager.recentLogLevel)
	assert.Equal(t, infoLogLevelCode, logManager.recentLogLevelCode)
	assert.Equal(t, expectedTestMessage, logManager.logs[0].logMessage)
	logManager.clearLogs()
}

func TestLoggerPrintln(t *testing.T) {
	Logger.Println(testMessage)
	assert.Equal(t, infoLogLevel, logManager.recentLogLevel)
	assert.Equal(t, infoLogLevelCode, logManager.recentLogLevelCode)
	assert.Equal(t, expectedTestMessage, logManager.logs[0].logMessage)
	logManager.clearLogs()
}

func TestLoggerPanicf(t *testing.T) {
	panicTestFunc := func() {
		Logger.Panicf("[%s]", testMessage)
	}
	assert.Panics(t, panicTestFunc)
	assert.Equal(t, errorLogLevel, logManager.recentLogLevel)
	assert.Equal(t, errorLogLevelCode, logManager.recentLogLevelCode)
	assert.Equal(t, formattedTestMessage, logManager.logs[0].logMessage)
	logManager.clearLogs()
}

func TestLoggerPanic(t *testing.T) {
	panicTestFunc := func() {
		Logger.Panic(testMessage)
	}
	assert.Panics(t, panicTestFunc)
	assert.Equal(t, errorLogLevel, logManager.recentLogLevel)
	assert.Equal(t, errorLogLevelCode, logManager.recentLogLevelCode)
	assert.Equal(t, expectedTestMessage, logManager.logs[0].logMessage)
	logManager.clearLogs()
}

func TestLoggerPanicln(t *testing.T) {
	panicTestFunc := func() {
		Logger.Panicln(testMessage)
	}
	assert.Panics(t, panicTestFunc)
	assert.Equal(t, errorLogLevel, logManager.recentLogLevel)
	assert.Equal(t, errorLogLevelCode, logManager.recentLogLevelCode)
	assert.Equal(t, expectedTestMessage, logManager.logs[0].logMessage)
	logManager.clearLogs()
}

func TestLogManagerWrite(t *testing.T) {
	p1, err1 := json.Marshal("testMessage")
	p2, err2 := json.Marshal("anotherTestMessage")
	if err1 != nil || err2 != nil {
		log.Println(err1, err2)
	}
	logManager.Write(p1)
	logManager.Write(p2)

	testMonitoredLogSetCorrectly(t, logManager.logs[0], "\"testMessage\"")
	testMonitoredLogSetCorrectly(t, logManager.logs[1], "\"anotherTestMessage\"")
}

func testMonitoredLogSetCorrectly(t *testing.T, m *monitoringLog, expectedMessage string) {
	assert.Equal(t, expectedMessage, m.logMessage)

	now := utils.GetTimestamp()
	assert.True(t, now >= m.logTimestamp)
}
