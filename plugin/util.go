package plugin

import (
	"github.com/satori/go.uuid"
	"strings"
	"github.com/aws/aws-lambda-go/lambdacontext"
	"os"
	"github.com/shirou/gopsutil/process"
	"fmt"
	"time"
)

func GenerateNewId() string {
	return uuid.NewV4().String()
}

//AppId starts after ] in logstreamname
func GetAppIdFromStreamName(logStreamName string) string {
	s := strings.Split(logStreamName, "]")
	if len(s) > 1 {
		return s[1]
	} else {
		return ""
	}
}

func GetApplicationVersion() string {
	return lambdacontext.FunctionVersion
}

func GetApplicationProfile() string {
	p := os.Getenv(ThundraApplicationProfile)
	if p == "" {
		p = DefaultProfile
	}
	return p
}

func GetApplicationName() string {
	return lambdacontext.FunctionName
}

func GetApplicationType() string {
	return ApplicationType
}

func GetThisProcess() *process.Process {
	pid := os.Getpid()
	p, err := process.NewProcess(int32(pid))
	if err != nil {
		fmt.Println(err)
		return nil
	}
	return p
}

//Returns current Unix timestamp in msec
func GetTimestamp() int64 {
	return time.Now().UnixNano() / (int64(time.Millisecond) / int64(time.Nanosecond))
}
