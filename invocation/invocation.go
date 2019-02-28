package invocation

import (
	"context"
	"encoding/json"

	opentracing "github.com/opentracing/opentracing-go"
	"github.com/thundra-io/thundra-lambda-agent-go/plugin"
	"github.com/thundra-io/thundra-lambda-agent-go/utils"
)

var invocationCount uint32

type invocationPlugin struct {
	data     *invocationData
	rootSpan opentracing.Span
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
	return &invocationPlugin{
		data: &invocationData{},
	}
}

func (ip *invocationPlugin) IsEnabled() bool {
	return true
}

func (ip *invocationPlugin) Order() uint8 {
	return pluginOrder
}

func (ip *invocationPlugin) BeforeExecution(ctx context.Context, request json.RawMessage) context.Context {
	ip.rootSpan = opentracing.SpanFromContext(ctx)
	ip.data = &invocationData{
		startTimestamp: utils.GetTimestamp(),
	}
	return ctx
}

func (ip *invocationPlugin) AfterExecution(ctx context.Context, request json.RawMessage, response interface{}, err interface{}) []plugin.MonitoringDataWrapper {
	ip.data.finishTimestamp = utils.GetTimestamp()
	ip.data.duration = ip.data.finishTimestamp - ip.data.startTimestamp

	if err != nil {
		ip.data.erroneous = true
		ip.data.errorMessage = utils.GetErrorMessage(err)
		ip.data.errorType = utils.GetErrorType(err)
		ip.data.errorCode = defaultErrorCode
	}

	ip.data.coldStart = isColdStarted()
	ip.data.timeout = utils.IsTimeout(err)

	data := ip.prepareData(ctx)

	ip.Reset()

	return []plugin.MonitoringDataWrapper{plugin.WrapMonitoringData(data, "Invocation")}
}

func (ip *invocationPlugin) Reset() {
	ClearTags()
}

// isColdStarted returns if the lambda instance is cold started. Cold Start only happens on the first invocationPlugin.
func isColdStarted() (coldStart bool) {
	if invocationCount++; invocationCount == 1 {
		coldStart = true
	}
	return coldStart
}
