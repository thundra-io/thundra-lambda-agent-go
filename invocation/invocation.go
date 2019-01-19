package invocation

import (
	"context"
	"encoding/json"
	"sync"

	"github.com/thundra-io/thundra-lambda-agent-go/plugin"
)

var invocationCount uint32

type invocationPlugin struct {
	data *invocationData
}

type invocationData struct {
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

// New initializes and returns a new invocationPlugin object.
func New() *invocationPlugin {
	ip := new(invocationPlugin)
	ip.data = new(invocationData)
	return ip
}

func (ip *invocationPlugin) IsEnabled() bool {
	return true
}

func (ip *invocationPlugin) BeforeExecution(ctx context.Context, request json.RawMessage, wg *sync.WaitGroup) {
	ip.data = new(invocationData)
	ip.data.startTimestamp = plugin.GetTimestamp()
	wg.Done()
}

func (ip *invocationPlugin) AfterExecution(ctx context.Context, request json.RawMessage, response interface{}, err interface{}) []plugin.MonitoringDataWrapper {
	ip.data.finishTimestamp = plugin.GetTimestamp()
	ip.data.duration = ip.data.finishTimestamp - ip.data.startTimestamp

	if err != nil {
		ip.data.erroneous = true
		ip.data.errorMessage = plugin.GetErrorMessage(err)
		ip.data.errorType = plugin.GetErrorType(err)
		ip.data.errorCode = defaultErrorCode
	}

	ip.data.coldStart = isColdStarted()
	ip.data.timeout = plugin.IsTimeout(err)

	data := ip.prepareData(ctx)
	ip.data = nil
	var invocationArr []plugin.MonitoringDataWrapper
	invocationArr = append(invocationArr, plugin.WrapMonitoringData(data, "Invocation"))
	return invocationArr
}

// isColdStarted returns if the lambda instance is cold started. Cold Start only happens on the first invocationPlugin.
func isColdStarted() (coldStart bool) {
	if invocationCount++; invocationCount == 1 {
		coldStart = true
	}
	return coldStart
}
