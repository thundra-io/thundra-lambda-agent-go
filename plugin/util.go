package plugin

import (
	"fmt"
	"os"
	"time"

	"github.com/satori/go.uuid"
	"github.com/shirou/gopsutil/process"
	"reflect"
)

// GenerateNewId generates new uuid.
func GenerateNewId() string {
	return uuid.NewV4().String()
}

// Generate2NewId generates 2 new uuids.
func Generate2NewId() (string, string) {
	return GenerateNewId(), GenerateNewId()
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
