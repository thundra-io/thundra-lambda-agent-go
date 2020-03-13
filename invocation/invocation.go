package invocation

import (
	"context"
	"encoding/json"

	opentracing "github.com/opentracing/opentracing-go"
	"github.com/thundra-io/thundra-lambda-agent-go/constants"
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
	startTime, ctx := plugin.StartTimeFromContext(ctx)
	ip.rootSpan = opentracing.SpanFromContext(ctx)
	ip.data = &invocationData{
		startTimestamp: startTime,
	}

	setInvocationTriggerTags(ctx, request)
	if GetAgentTag(constants.SpanTags["TRIGGER_CLASS_NAME"]) != nil {
		triggerClassName, ok := GetAgentTag(constants.SpanTags["TRIGGER_CLASS_NAME"]).(string)
		if ok {
			plugin.TriggerClassName = triggerClassName
		}
	}
	return ctx
}

func (ip *invocationPlugin) AfterExecution(ctx context.Context, request json.RawMessage, response interface{}, err interface{}) ([]plugin.MonitoringDataWrapper, context.Context) {
	finishTime, ctx := plugin.EndTimeFromContext(ctx)
	ip.data.finishTimestamp = finishTime
	ip.data.duration = ip.data.finishTimestamp - ip.data.startTimestamp

	if userError != nil {
		ip.data.erroneous = true
		ip.data.errorMessage = utils.GetErrorMessage(userError)
		ip.data.errorType = utils.GetErrorType(userError)
		ip.data.errorCode = defaultErrorCode
		utils.SetSpanError(ip.rootSpan, userError)
	}

	if err != nil {
		ip.data.erroneous = true
		ip.data.errorMessage = utils.GetErrorMessage(err)
		ip.data.errorType = utils.GetErrorType(err)
		ip.data.errorCode = defaultErrorCode
	}

	ip.data.coldStart = isColdStarted()
	ip.data.timeout = utils.IsTimeout(err)

	data := ip.prepareData(ctx)

	if response != nil {
		responseInterface, ok := (response).(*interface{})
		if ok {
			statusCode, err := utils.GetStatusCode(responseInterface)
			if err == nil {
				SetTag(constants.HTTPTags["STATUS"], statusCode)
			}
		}
	}

	ip.Reset()

	return []plugin.MonitoringDataWrapper{plugin.WrapMonitoringData(data, "Invocation")}, ctx
}

func (ip *invocationPlugin) Reset() {
	Clear()
	clearTraceLinks()
}

// isColdStarted returns if the lambda instance is cold started. Cold Start only happens on the first invocationPlugin.
func isColdStarted() (coldStart bool) {
	if invocationCount++; invocationCount == 1 {
		coldStart = true
	}
	return coldStart
}
