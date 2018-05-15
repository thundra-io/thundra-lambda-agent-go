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

// GetAppId returns application id.
func GetAppId() string {
	return getAppIdFromStreamName(lambdacontext.LogStreamName)
}

// getAppIdFromStreamName returns application id. AppId starts after ']' in logstreamname.
func getAppIdFromStreamName(logStreamName string) string {
	s := strings.Split(logStreamName, "]")
	if len(s) > 1 {
		return s[1]
	} else {
		return ""
	}
}

// GetApplicationVersion returns function version.
func GetApplicationVersion() string {
	return lambdacontext.FunctionVersion
}

// GetApplicationProfile returns profile.
func GetApplicationProfile() string {
	p := os.Getenv(ThundraApplicationProfile)
	if p == "" {
		p = DefaultProfile
	}
	return p
}

// GetApplicationName return function name.
func GetApplicationName() string {
	return lambdacontext.FunctionName
}

// GetApplicationType returns programming language type, i.e. "go."
func GetApplicationType() string {
	return ApplicationType
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

// GetRegion returns AWS region's name
func GetRegion() string {
	return os.Getenv(AwsDefaultRegion)
}

// GetMemorySize returns configured memory limit for the current instance of the Lambda Function
func GetMemorySize() int {
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
	return err.(error).Error()
}
