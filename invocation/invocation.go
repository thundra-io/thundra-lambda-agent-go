package invocation

import (
	"context"
	"encoding/json"
	"sync"

	"github.com/thundra-io/thundra-lambda-agent-go/plugin"
)

var invocationCount uint32

// invocationSpan collects raw information related to invocation.
type invocationSpan struct {
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

type invocation struct {
	span *invocationSpan
}

// New initializes and returns a new invocation object.
func New() *invocation {
	i := new(invocation)
	i.span = new(invocationSpan)
	return i
}

func (i *invocation) BeforeExecution(ctx context.Context, request json.RawMessage, wg *sync.WaitGroup) {
	i.span = new(invocationSpan)
	i.span.startTimestamp = plugin.GetTimestamp()
	wg.Done()
}

func (i *invocation) AfterExecution(ctx context.Context, request json.RawMessage, response interface{}, err interface{}) []plugin.MonitoringDataWrapper {
	i.span.finishTimestamp = plugin.GetTimestamp()
	i.span.duration = i.span.finishTimestamp - i.span.startTimestamp

	if err != nil {
		i.span.erroneous = true
		i.span.errorMessage = plugin.GetErrorMessage(err)
		i.span.errorType = plugin.GetErrorType(err)
		i.span.errorCode = defaultErrorCode
	}

	i.span.coldStart = isColdStarted()
	i.span.timeout = plugin.IsTimeout(err)

	data := i.prepareData(ctx)
	i.span = nil

	var invocationArr []plugin.MonitoringDataWrapper
	invocationArr = append(invocationArr, plugin.WrapMonitoringData(data, invocationType))
	return invocationArr
}

func (i *invocation) OnPanic(ctx context.Context, request json.RawMessage, err interface{}, stackTrace []byte) []plugin.MonitoringDataWrapper {
	i.span.finishTimestamp = plugin.GetTimestamp()
	i.span.duration = i.span.finishTimestamp - i.span.startTimestamp
	i.span.erroneous = true
	i.span.errorMessage = plugin.GetErrorMessage(err)
	i.span.errorType = plugin.GetErrorType(err)
	i.span.errorCode = defaultErrorCode
	i.span.coldStart = isColdStarted()

	// since it is panicked it could not be timed out
	i.span.timeout = false

	data := i.prepareData(ctx)
	i.span = nil

	var invocationArr []plugin.MonitoringDataWrapper
	invocationArr = append(invocationArr, plugin.WrapMonitoringData(data, invocationType))
	return invocationArr
}

// isColdStarted returns if the lambda instance is cold started. Cold Start only happens on the first invocation.
func isColdStarted() (coldStart bool) {
	if invocationCount++; invocationCount == 1 {
		coldStart = true
	}
	return coldStart
}