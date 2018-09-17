package invocation

import (
	"context"
	"encoding/json"
	"sync"

	"github.com/thundra-io/thundra-lambda-agent-go/plugin"
)

var invocationCount uint32

// New initializes and returns a new invocation object.
func New() *invocation {
	return &invocation{
		Type:                      invocationType,
		AgentVersion:              plugin.AgentVersion,
		DataModelVersion:          plugin.DataModelVersion,
		ApplicationId:             plugin.ApplicationId,
		ApplicationDomainName:     plugin.ApplicationDomainName,
		ApplicationClassName:      plugin.ApplicationClassName,
		ApplicationName:           plugin.FunctionName,
		ApplicationVersion:        plugin.ApplicationVersion,
		ApplicationStage:          plugin.ApplicationStage,
		ApplicationRuntime:        plugin.ApplicationRuntime,
		ApplicationRuntimeVersion: plugin.ApplicationRuntimeVersion,
		ApplicationTags:           map[string]interface{}{}, // empty object

		FunctionPlatform: functionPlatform,
		FunctionName:     plugin.FunctionName,
		FunctionRegion:   plugin.FunctionRegion,
		Tags:             map[string]interface{}{}, // empty object
	}
}

func (i *invocation) BeforeExecution(ctx context.Context, request json.RawMessage, wg *sync.WaitGroup) {
	i.Id = plugin.GenerateNewId()
	i.TraceId = plugin.TraceId
	i.TransactionId = plugin.TransactionId
	i.StartTimestamp = plugin.GetTimestamp()
	wg.Done()
}

func (i *invocation) AfterExecution(ctx context.Context, request json.RawMessage, response interface{}, err interface{}) ([]interface{}, string) {
	i.FinishTimestamp = plugin.GetTimestamp()
	i.Duration = i.FinishTimestamp - i.StartTimestamp

	if err != nil {
		i.Erroneous = true
		i.ErrorMessage = plugin.GetErrorMessage(err)
		i.ErrorType = plugin.GetErrorType(err)
		i.ErrorCode = defaultErrorCode
	}

	i.ColdStart = isColdStarted()
	i.Timeout = isTimeout(err)

	var invocationArr []interface{}
	invocationArr = append(invocationArr, i)
	return invocationArr, invocationType
}

func (i *invocation) OnPanic(ctx context.Context, request json.RawMessage, err interface{}, stackTrace []byte) ([]interface{}, string) {
	i.FinishTimestamp = plugin.GetTimestamp()
	i.Duration = i.FinishTimestamp - i.StartTimestamp
	i.Erroneous = true
	i.ErrorMessage = plugin.GetErrorMessage(err)
	i.ErrorType = plugin.GetErrorType(err)
	i.ErrorCode = defaultErrorCode
	i.ColdStart = isColdStarted()

	// since it is panicked it could not be timed out
	i.Timeout = false

	var invocationArr []interface{}
	invocationArr = append(invocationArr, i)
	return invocationArr, invocationType
}

// isColdStarted returns if the lambda instance is cold started. Cold Start only happens on the first invocation.
func isColdStarted() (coldStart bool) {
	if invocationCount += 1; invocationCount == 1 {
		coldStart = true
	}
	return coldStart
}

// isTimeout returns if the lambda invocation is timed out.
func isTimeout(err interface{}) bool {
	if err == nil {
		return false
	}
	if plugin.GetErrorType(err) == "timeoutError" {
		return true
	}
	return false
}
