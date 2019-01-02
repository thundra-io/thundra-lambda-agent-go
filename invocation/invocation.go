package invocation

import (
	"context"
	"encoding/json"
	"sync"

	"github.com/thundra-io/thundra-lambda-agent-go/plugin"
)

var invocationCount uint32

type invocation struct {
	startTimestamp  int64
	finishTimestamp int64
	duration        int64
	erroneous       bool
	errorMessage    string
	errorType       string
	errorCode       string
	coldStart       bool
	timeout         bool
}

// New initializes and returns a new invocation object.
func New() *invocation {
	i := new(invocation)
	return i
}

func (i *invocation) BeforeExecution(ctx context.Context, request json.RawMessage, wg *sync.WaitGroup) {
	i.startTimestamp = plugin.GetTimestamp()
	wg.Done()
}

func (i *invocation) AfterExecution(ctx context.Context, request json.RawMessage, response interface{}, err interface{}) []plugin.MonitoringDataWrapper {
	i.finishTimestamp = plugin.GetTimestamp()
	i.duration = i.finishTimestamp - i.startTimestamp

	if err != nil {
		i.erroneous = true
		i.errorMessage = plugin.GetErrorMessage(err)
		i.errorType = plugin.GetErrorType(err)
		i.errorCode = defaultErrorCode
	}

	i.coldStart = isColdStarted()
	i.timeout = plugin.IsTimeout(err)

	data := i.prepareData(ctx)
	i = nil

	var invocationArr []plugin.MonitoringDataWrapper
	invocationArr = append(invocationArr, plugin.WrapMonitoringData(data, "Invocation"))
	return invocationArr
}

func (i *invocation) OnPanic(ctx context.Context, request json.RawMessage, err interface{}, stackTrace []byte) []plugin.MonitoringDataWrapper {
	i.finishTimestamp = plugin.GetTimestamp()
	i.duration = i.finishTimestamp - i.startTimestamp
	i.erroneous = true
	i.errorMessage = plugin.GetErrorMessage(err)
	i.errorType = plugin.GetErrorType(err)
	i.errorCode = defaultErrorCode
	i.coldStart = isColdStarted()

	// since it is panicked it could not be timed out
	i.timeout = false

	data := i.prepareData(ctx)
	i = nil

	var invocationArr []plugin.MonitoringDataWrapper
	invocationArr = append(invocationArr, plugin.WrapMonitoringData(data, "Invocation"))
	return invocationArr
}

// isColdStarted returns if the lambda instance is cold started. Cold Start only happens on the first invocation.
func isColdStarted() (coldStart bool) {
	if invocationCount++; invocationCount == 1 {
		coldStart = true
	}
	return coldStart
}