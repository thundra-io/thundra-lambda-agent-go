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
		ApplicationName:    plugin.ApplicationName,
		ApplicationId:      plugin.ApplicationId,
		ApplicationVersion: plugin.ApplicationVersion,
		ApplicationProfile: plugin.ApplicationProfile,
		ApplicationType:    plugin.ApplicationType,
		Region:             plugin.Region,
		MemorySize:         plugin.MemorySize,
	}
}

func (i *invocation) BeforeExecution(ctx context.Context, request json.RawMessage, wg *sync.WaitGroup) {
	i.Id = plugin.GenerateNewId()
	i.TransactionId = plugin.TransactionId
	i.StartTimestamp = plugin.GetTimestamp()
	wg.Done()
}

func (i *invocation) AfterExecution(ctx context.Context, request json.RawMessage, response interface{}, err interface{}) ([]interface{}, string) {
	i.EndTimestamp = plugin.GetTimestamp()
	i.Duration = i.EndTimestamp - i.StartTimestamp

	if err != nil {
		i.Erroneous = true
		i.ErrorMessage = plugin.GetErrorMessage(err)
		i.ErrorType = plugin.GetErrorType(err)
	}

	i.ColdStart = isColdStarted()
	i.Timeout = isTimeout()

	var invocationArr []interface{}
	invocationArr = append(invocationArr, i)
	return invocationArr, invocationDataType
}

func (i *invocation) OnPanic(ctx context.Context, request json.RawMessage, err interface{}, stackTrace []byte) ([]interface{}, string) {
	i.EndTimestamp = plugin.GetTimestamp()
	i.Duration = i.EndTimestamp - i.StartTimestamp
	i.Erroneous = true
	i.ErrorMessage = plugin.GetErrorMessage(err)
	i.ErrorType = plugin.GetErrorType(err)
	i.ColdStart = isColdStarted()
	i.Timeout = isTimeout()

	var invocationArr []interface{}
	invocationArr = append(invocationArr, i)
	return invocationArr, invocationDataType
}

// isColdStarted returns if the lambda instance is cold started. Cold Start only happens on the first invocation.
func isColdStarted() (coldStart bool) {
	if invocationCount += 1; invocationCount == 1 {
		coldStart = true
	}
	return coldStart
}

// TODO !NOT IMPLEMENTED YET!
// isTimeout returns if the lambda invocation is timed out.
func isTimeout() (timeout bool) {
	return timeout
}
