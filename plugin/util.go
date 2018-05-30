package plugin

import (
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/aws/aws-lambda-go/lambdacontext"
	"github.com/satori/go.uuid"
	"github.com/shirou/gopsutil/process"
	"reflect"
)

// GenerateNewId generates new uuid.
func GenerateNewId() string {
	return uuid.NewV4().String()
}

func GenerateNewTransactionId() {
	TransactionId = GenerateNewId()
}

// getAppId returns application id.
func getAppId() string {
	return getAppIdFromStreamName(lambdacontext.LogStreamName)
}

// getAppIdFromStreamName returns application id. AppId starts after ']' in logstreamname.
func getAppIdFromStreamName(logStreamName string) string {
	s := strings.Split(logStreamName, "]")
	if len(s) > 1 {
		return s[1]
	}
	return ""
}

// getApplicationVersion returns function version.
func getApplicationVersion() string {
	return lambdacontext.FunctionVersion
}

// getApplicationProfile returns profile.
func getApplicationProfile() string {
	p := os.Getenv(ThundraApplicationProfile)
	if p == "" {
		p = DefaultProfile
	}
	return p
}

// getApplicationName returns function name.
func getApplicationName() string {
	return lambdacontext.FunctionName
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

// getRegion returns AWS region's name
func getRegion() string {
	return os.Getenv(AwsDefaultRegion)
}

// getMemorySize returns configured memory limit for the current instance of the Lambda Function
func getMemorySize() int {
	return lambdacontext.MemoryLimitInMB
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
