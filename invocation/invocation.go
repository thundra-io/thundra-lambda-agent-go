package invocation

import (
	"context"
	"encoding/json"
	"sync"

	"github.com/thundra-io/thundra-lambda-agent-go/plugin"
)

var invocationCount uint32

type invocation struct {
	Id                 string `json:"id"`
	TransactionId      string `json:"transactionId"`
	ApplicationName    string `json:"applicationName"`
	ApplicationId      string `json:"applicationId"`
	ApplicationVersion string `json:"applicationVersion"`
	ApplicationProfile string `json:"applicationProfile"`
	ApplicationType    string `json:"applicationType"`

	Duration       int64  `json:"duration"`       // Invocation time in milliseconds
	StartTimestamp int64  `json:"startTimestamp"` // Invocation start time in UNIX Epoch milliseconds
	EndTimestamp   int64  `json:"endTimestamp"`   // Invocation end time in UNIX Epoch milliseconds
	Erroneous      bool   `json:"erroneous"`      // Shows if the invocation failed with an error
	ErrorType      string `json:"errorType"`      // Type of the thrown error
	ErrorMessage   string `json:"errorMessage"`   // Message of the thrown error
	ColdStart      bool   `json:"coldStart"`      // Shows if the invocation is cold started
	Timeout        bool   `json:"timeout"`        // Shows if the invocation is timed out

	Region     string `json:"region"`     // Name of the AWS region
	MemorySize int    `json:"memorySize"` // Memory Size of the function in MB
}

// NewInvocation initializes and returns a new invocation object
func NewInvocation() *invocation {
	return &invocation{
		ApplicationName:    plugin.GetApplicationName(),
		ApplicationId:      plugin.GetAppId(),
		ApplicationVersion: plugin.GetApplicationVersion(),
		ApplicationProfile: plugin.GetApplicationProfile(),
		ApplicationType:    plugin.GetApplicationType(),
		Region:             plugin.GetRegion(),
		MemorySize:         plugin.GetMemorySize(),
	}
}

func (i *invocation) BeforeExecution(ctx context.Context, request json.RawMessage, transactionId string, wg *sync.WaitGroup) {
	i.Id = plugin.GenerateNewId()
	i.TransactionId = transactionId
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
	i.ErrorMessage = err.(error).Error()
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

// !NOT IMPLEMENTED YET!
// isTimeout returns if the lambda invocation is timed out.
func isTimeout() (timeout bool) {
	return timeout
}
