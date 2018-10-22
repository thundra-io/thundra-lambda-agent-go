package plugin

import (
	"context"
	"os"
	"strings"

	"github.com/aws/aws-lambda-go/lambdacontext"
)

var FunctionName string
var ApplicationId string
var ApplicationVersion string
var ApplicationStage string
var FunctionRegion string
var MemoryLimit int
var LogGroupName string
var LogStreamName string
var FunctionARN string

var TraceId string
var TransactionId string
var SpanId string

var ApiKey string

func init() {
	FunctionName = getFunctionName()
	ApplicationId = getAppId()
	ApplicationVersion = getApplicationVersion()
	ApplicationStage = getApplicationStage()
	FunctionRegion = getFunctionRegion()
	MemoryLimit = getMemoryLimit()
	LogGroupName = getLogGroupName()
	LogStreamName = getLogStreamName()
}

// getFunctionName returns function name.
func getFunctionName() string {
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

// getApplicationStage returns profile.
func getApplicationStage() string {
	p := os.Getenv(ThundraApplicationProfile)
	if p == "" {
		p = DefaultProfile
	}
	return p
}

// getFunctionRegion returns AWS region's name
func getFunctionRegion() string {
	return os.Getenv(AwsDefaultRegion)
}

// getMemoryLimit returns configured memory limit for the current instance of the Lambda Function
func getMemoryLimit() int {
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
func GetAwsRequestID(ctx context.Context) string {
	lc, ok := lambdacontext.FromContext(ctx)
	if !ok {
		// lambdaContext is not set
		return ""
	}
	return lc.AwsRequestID
}
