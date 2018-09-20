package trace

import (
	"context"
	"encoding/json"

	"github.com/thundra-io/thundra-lambda-agent-go/plugin"
)

type traceData struct {
	//Base fields
	Id                        string                 `json:"id"`
	Type                      string                 `json:"type"`
	AgentVersion              string                 `json:"agentVersion"`
	DataModelVersion          string                 `json:"dataModelVersion"`
	ApplicationId             string                 `json:"applicationId"`
	ApplicationDomainName     string                 `json:"applicationDomainName"`
	ApplicationClassName      string                 `json:"applicationClassName"`
	ApplicationName           string                 `json:"applicationName"`
	ApplicationVersion        string                 `json:"applicationVersion"`
	ApplicationStage          string                 `json:"applicationStage"`
	ApplicationRuntime        string                 `json:"applicationRuntime"`
	ApplicationRuntimeVersion string                 `json:"applicationRuntimeVersion"`
	ApplicationTags           map[string]interface{} `json:"applicationTags"`

	//Trace fields
	RootSpanId      string                 `json:"rootSpanId"`
	StartTimestamp  int64                  `json:"startTimestamp"`
	FinishTimestamp int64                  `json:"finishTimestamp"`
	Duration        int64                  `json:"duration"`
	Tags            map[string]interface{} `json:"tags"`
}

func (tr *trace) prepareTraceData(ctx context.Context, request json.RawMessage, response interface{}) traceData {
	tags := tr.prepareTraceTags(ctx, request, response)
	return traceData{
		Id:                        plugin.TraceId,
		Type:                      traceType,
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
		ApplicationTags:           map[string]interface{}{},

		RootSpanId:      tr.span.rootSpanId,
		StartTimestamp:  tr.span.startTime,
		FinishTimestamp: tr.span.finishTime,
		Duration:        tr.span.duration,
		Tags:            tags,
	}
}

func (tr *trace) prepareTraceTags(ctx context.Context, request json.RawMessage, response interface{}) map[string]interface{} {
	tags := map[string]interface{}{}
	tags[plugin.AwsLambdaInvocationRequestId] = plugin.GetAwsRequestID(ctx)

	// If the agent's user doesn't want to send their request and response data, hide them.
	if !shouldHideRequest() {
		tags[plugin.AwsLambdaInvocationRequest] = string(request)
	}
	if !shouldHideResponse() {
		tags[plugin.AwsLambdaInvocationResponse] = response
	}

	tags[plugin.AwsLambdaARN] = plugin.GetInvokedFunctionArn(ctx)
	tags[plugin.AwsLambdaLogGroupName] = plugin.LogGroupName
	tags[plugin.AwsLambdaLogStreamName] = plugin.LogStreamName
	tags[plugin.AwsLambdaMemoryLimit] = plugin.MemoryLimit
	tags[plugin.AwsLambdaName] = plugin.FunctionName
	tags[plugin.AwsRegion] = plugin.FunctionRegion
	tags[plugin.AwsLambdaInvocationTimeout] = tr.span.timeout

	// If this is the first invocation, it is a cold start
	if invocationCount == 1 {
		tags[plugin.AwsLambdaInvocationColdStart] = true
	} else {
		tags[plugin.AwsLambdaInvocationColdStart] = false
	}

	if tr.span.panicInfo != nil {
		tags[plugin.AwsError] = true
		tags[plugin.AwsErrorKind] = tr.span.panicInfo.Kind
		tags[plugin.AwsErrorMessage] = tr.span.panicInfo.Message
		tags[plugin.AwsErrorStack] = tr.span.panicInfo.Stack
	} else if tr.span.errorInfo != nil {
		tags[plugin.AwsError] = true
		tags[plugin.AwsErrorKind] = tr.span.errorInfo.Kind
		tags[plugin.AwsErrorMessage] = tr.span.errorInfo.Message
	}
	return tags
}

type spanData struct {
	//Base fields
	Id                        string                 `json:"id"`
	Type                      string                 `json:"type"`
	AgentVersion              string                 `json:"agentVersion"`
	DataModelVersion          string                 `json:"dataModelVersion"`
	ApplicationId             string                 `json:"applicationId"`
	ApplicationDomainName     string                 `json:"applicationDomainName"`
	ApplicationClassName      string                 `json:"applicationClassName"`
	ApplicationName           string                 `json:"applicationName"`
	ApplicationVersion        string                 `json:"applicationVersion"`
	ApplicationStage          string                 `json:"applicationStage"`
	ApplicationRuntime        string                 `json:"applicationRuntime"`
	ApplicationRuntimeVersion string                 `json:"applicationRuntimeVersion"`
	ApplicationTags           map[string]interface{} `json:"applicationTags"`

	TraceId        string `json:"traceId"`
	TracnsactionId string `json:"transactionId"`
	ParentSpanId   string `json:"parentSpanId"`

	SpanOrder       int64                  `json:"spanOrder"`
	DomainName      string                 `json:"domainName"`
	ClassName       string                 `json:"className"`
	ServiceName     string                 `json:"serviceName"`
	OperationName   string                 `json:"operationName"`
	StartTimestamp  int64                  `json:"startTimestamp"`
	FinishTimestamp int64                  `json:"finishTimestamp"`
	Duration        int64                  `json:"duration"`
	Tags            map[string]interface{} `json:"tags"`
	Logs            map[string]spanLog     `json:"logs"`
}

type spanLog struct {
	Name      string      `json:"name"`
	Value     interface{} `json:"value"`
	Timestamp int64       `json:"timestamp"`
}

func (tr *trace) prepareSpanData(ctx context.Context, request json.RawMessage, response interface{}) spanData {
	tags := tr.prepareSpanTags(ctx, request, response)
	return spanData{
		Id:                        plugin.GenerateNewId(),
		Type:                      spanType,
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
		ApplicationTags:           map[string]interface{}{},

		TraceId:        plugin.TraceId,
		TracnsactionId: plugin.TransactionId,
		//ParentSpanId:    Only root span exist right now

		DomainName:    plugin.ApplicationDomainName,
		ClassName:     plugin.ApplicationClassName,
		ServiceName:   plugin.FunctionName, //TODO implement it with Opentracing
		OperationName: plugin.FunctionName, //TODO implement it with Opentracing

		StartTimestamp:  tr.span.startTime,
		FinishTimestamp: tr.span.finishTime,
		Duration:        tr.span.duration,
		Tags:            tags,
	}
}

func (tr *trace) prepareSpanTags(ctx context.Context, request json.RawMessage, response interface{}) map[string]interface{} {
	tags := map[string]interface{}{}
	// If the agent's user doesn't want to send their request and response data, hide them.
	if !shouldHideRequest() {
		tags[plugin.AwsLambdaInvocationRequest] = string(request)
	}
	if !shouldHideResponse() {
		tags[plugin.AwsLambdaInvocationResponse] = response
	}
	return tags
}
