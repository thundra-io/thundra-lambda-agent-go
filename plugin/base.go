package plugin

import (
	"github.com/aws/aws-lambda-go/lambdacontext"
	"os"
	"strings"
	"context"
)

var ApplicationName string
var ApplicationId string
var ApplicationVersion string
var ApplicationProfile string
var Region string
var MemorySize int
var LogGroupName string
var LogStreamName string
var FunctionARN string

var TransactionId string
var ContextId string

func init() {
	ApplicationName = getApplicationName()
	ApplicationId = getAppId()
	ApplicationVersion = getApplicationVersion()
	ApplicationProfile = getApplicationProfile()
	Region = getRegion()
	MemorySize = getMemorySize()
	LogGroupName = getLogGroupName()
	LogStreamName = getLogStreamName()
}

// getApplicationName returns function name.
func getApplicationName() string {
	return lambdacontext.FunctionName
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

// getRegion returns AWS region's name
func getRegion() string {
	return os.Getenv(AwsDefaultRegion)
}

// getMemorySize returns configured memory limit for the current instance of the Lambda Function
func getMemorySize() int {
	return lambdacontext.MemoryLimitInMB
}

func getLogGroupName() string {
	return lambdacontext.LogGroupName
}

func getLogStreamName() string {
	return lambdacontext.LogStreamName
}

// GetFromContext returns InvokedFunctionArn and AwsRequestID if available.
func GetInvokedFunctionArn(ctx context.Context) string {
	lc, ok := lambdacontext.FromContext(ctx)
	if !ok {
		// lambdaContext is not set
		return ""
	}
	return lc.InvokedFunctionArn
}

// GetFromContext returns InvokedFunctionArn and AwsRequestID if available.
func GetAwsRequestID(ctx context.Context) (string) {
	lc, ok := lambdacontext.FromContext(ctx)
	if !ok {
		// lambdaContext is not set
		return ""
	}
	return lc.AwsRequestID
}
