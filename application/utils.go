package application

import (
	"context"
	"github.com/thundra-io/thundra-lambda-agent-go/constants"
	"os"
	"strings"

	"github.com/aws/aws-lambda-go/lambdacontext"
)

var ApplicationName string
var ApplicationID string
var ApplicationDomainName string
var ApplicationClassName string
var ApplicationVersion string
var ApplicationStage string
var FunctionRegion string
var MemoryLimit int
var LogGroupName string
var LogStreamName string
var FunctionARN string


func init() {
	ApplicationName = getApplicationName()
	ApplicationID = getAppID()
	ApplicationDomainName = getApplicationDomainName()
	ApplicationClassName = getApplicationClassName()
	ApplicationVersion = getApplicationVersion()
	ApplicationStage = getApplicationStage()
	FunctionRegion = getFunctionRegion()
	MemoryLimit = getMemoryLimit()
	LogGroupName = getLogGroupName()
	LogStreamName = getLogStreamName()
}

// getApplicationDomainName returns application domain name
func getApplicationDomainName() string {
	v := os.Getenv(constants.ApplicationDomainProp)
	if v != "" {
		return v
	}
	return constants.AwsLambdaApplicationDomain
}

// getApplicationClassName returns application class name
func getApplicationClassName() string {
	v := os.Getenv(constants.ApplicationClassProp)
	if v != "" {
		return v
	}
	return constants.AwsLambdaApplicationClass
}

// getApplicationName returns application name
func getApplicationName() string {
	v := os.Getenv(constants.ApplicationNameProp)
	if v != "" {
		return v
	}
	return lambdacontext.FunctionName
}

// getAppID returns application id
func getAppID() string {
	v := os.Getenv(constants.ApplicationIDProp)
	if v != "" {
		return v
	}
	return getAppIDFromStreamName(lambdacontext.LogStreamName)
}

// getAppIDFromStreamName returns application id. AppId starts after ']' in logstreamname
func getAppIDFromStreamName(logStreamName string) string {
	s := strings.Split(logStreamName, "]")
	if len(s) > 1 {
		return s[1]
	}
	return ""
}

// getApplicationVersion returns function version
func getApplicationVersion() string {
	v := os.Getenv(constants.ApplicationVersionProp)
	if v != "" {
		return v
	}
	return lambdacontext.FunctionVersion
}

// getApplicationStage returns profile
func getApplicationStage() string {
	v := os.Getenv(constants.ApplicationStageProp)
	if v != "" {
		return v
	}
	return os.Getenv(constants.ThundraApplicationStage)
}

// getFunctionRegion returns AWS region's name
func getFunctionRegion() string {
	return os.Getenv(constants.AwsDefaultRegion)
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

// GetInvokedFunctionArn returns InvokedFunctionArn if available.
func GetInvokedFunctionArn(ctx context.Context) string {
	lc, ok := lambdacontext.FromContext(ctx)
	if !ok {
		// lambdaContext is not set
		return ""
	}
	return lc.InvokedFunctionArn
}

// GetAwsRequestID returns AwsRequestID if available.
func GetAwsRequestID(ctx context.Context) string {
	lc, ok := lambdacontext.FromContext(ctx)
	if !ok {
		// lambdaContext is not set
		return ""
	}
	return lc.AwsRequestID
}
