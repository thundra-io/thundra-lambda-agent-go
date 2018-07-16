package plugin

import (
	"fmt"
	"os"
	"reflect"
	"strconv"
	"time"

	"github.com/satori/go.uuid"
	"github.com/shirou/gopsutil/process"
)

// GenerateNewId generates new uuid.
func GenerateNewId() string {
	return uuid.NewV4().String()
}

func GenerateNewTransactionId() {
	TransactionId = GenerateNewId()
}

func GenerateNewContextId() {
	ContextId = GenerateNewId()
}

// GetThisProcess returns process info about this process.
func GetThisProcess() *process.Process {
	pid := os.Getpid()
	p, err := process.NewProcess(int32(pid))
	if err != nil {
		fmt.Println(err)
		return nil
	}

	return p
}

// GetPid returns pid of this process.
func GetPid() string {
	pid := os.Getpid()
	return strconv.Itoa(pid)
}

// GetTimestamp returns current unix timestamp in msec.
func GetTimestamp() int64 {
	return time.Now().UnixNano() / (int64(time.Millisecond) / int64(time.Nanosecond))
}

// GetErrorType returns type of the error
func GetErrorType(err interface{}) string {
	errorType := reflect.TypeOf(err)
	if errorType.Kind() == reflect.Ptr {
		return errorType.Elem().Name()
	}
	return errorType.Name()
}

// GetErrorMessage returns error message
func GetErrorMessage(err interface{}) string {
	e, ok := err.(error)
	if !ok {
		return err.(string)
	}
	return e.Error()
}
