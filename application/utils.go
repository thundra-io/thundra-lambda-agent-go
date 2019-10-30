package application

import (
	"context"
	"fmt"
	"os"
	"strings"

	"github.com/thundra-io/thundra-lambda-agent-go/config"
	"github.com/thundra-io/thundra-lambda-agent-go/constants"

	"github.com/aws/aws-lambda-go/lambdacontext"
)

var ApplicationName string
var ApplicationInstanceID string
var FunctionName string
var ApplicationID string
var ApplicationDomainName string
var ApplicationClassName string
var ApplicationVersion string
var ApplicationStage string
var FunctionRegion string
var MemoryLimit int
var MemoryUsed int
var LogGroupName string
var LogStreamName string
var FunctionARN string
var ApplicationTags map[string]interface{}

func init() {
	ApplicationName = getApplicationName()
	ApplicationInstanceID = getApplicationInstanceID()
	ApplicationDomainName = getApplicationDomainName()
	ApplicationClassName = getApplicationClassName()
	ApplicationVersion = getApplicationVersion()
	ApplicationStage = getApplicationStage()
	FunctionRegion = getFunctionRegion()
	FunctionName = getFunctionName()
	MemoryLimit = getMemoryLimit()
	LogGroupName = getLogGroupName()
	LogStreamName = getLogStreamName()
	parseApplicationTags()
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

// getFunctionName returns the lambda function's name
func getFunctionName() string {
	return lambdacontext.FunctionName
}

func GetApplicationID(ctx context.Context) string {
	v := os.Getenv(constants.ApplicationIDProp)
	if v != "" {
		return v
	}

	arn := GetInvokedFunctionArn(ctx)
	accountNo := GetAwsAccountNo(arn)
	region := ""
	functionName := ""

	if len(FunctionRegion) > 0 {
		region = FunctionRegion
	} else if len(getFunctionRegion()) > 0 {
		region = getFunctionRegion()
	} else {
		region = "local"
	}

	if len(FunctionName) > 0 {
		functionName = FunctionName
	} else if len(getFunctionName()) > 0 {
		functionName = getFunctionName()
	} else {
		functionName = "lambda-app"
	}

	if config.SAMLocalDebugging {
		accountNo = "sam_local"
	} else if len(accountNo) == 0 {
		if len(config.APIKey) > 0 {
			accountNo = config.APIKey
		} else {
			accountNo = "guest"
		}
	}

	return fmt.Sprintf("aws:lambda:%s:%s:%s", region, accountNo, functionName)
}

func getApplicationInstanceID() string {
	logStreamName := getLogStreamName()
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

// GetAwsAccountNo returns Aws Account Number if available.
func GetAwsAccountNo(arn string) string {
	arnParts := strings.Split(arn, ":")
	if len(arnParts) > 4 {
		return arnParts[4]
	}
	return ""
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

// GetClientContext returns ClientContext object from context if available.
func GetClientContext(ctx context.Context) (lambdacontext.ClientContext, bool) {
	lc, ok := lambdacontext.FromContext(ctx)
	if !ok {
		// lambdaContext is not set
		return lambdacontext.ClientContext{}, ok
	}
	return lc.ClientContext, ok
}
