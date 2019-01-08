package trace

import (
	"context"
	"encoding/json"

	"github.com/thundra-io/thundra-lambda-agent-go/ttracer"

	"github.com/thundra-io/thundra-lambda-agent-go/plugin"
)

type traceDataModel struct {
	//Base fields
	ID                        string                 `json:"id"`
	Type                      string                 `json:"type"`
	AgentVersion              string                 `json:"agentVersion"`
	DataModelVersion          string                 `json:"dataModelVersion"`
	ApplicationID             string                 `json:"applicationId"`
	ApplicationDomainName     string                 `json:"applicationDomainName"`
	ApplicationClassName      string                 `json:"applicationClassName"`
	ApplicationName           string                 `json:"applicationName"`
	ApplicationVersion        string                 `json:"applicationVersion"`
	ApplicationStage          string                 `json:"applicationStage"`
	ApplicationRuntime        string                 `json:"applicationRuntime"`
	ApplicationRuntimeVersion string                 `json:"applicationRuntimeVersion"`
	ApplicationTags           map[string]interface{} `json:"applicationTags"`

	//Trace fields
	RootSpanID      string                 `json:"rootSpanId"`
	StartTimestamp  int64                  `json:"startTimestamp"`
	FinishTimestamp int64                  `json:"finishTimestamp"`
	Duration        int64                  `json:"duration"`
	Tags            map[string]interface{} `json:"tags"`
}

func (tr *trace) prepareTraceDataModel(ctx context.Context, request json.RawMessage, response interface{}) traceDataModel {
	tags := tr.prepareTraceTags(ctx, request, response)
	return traceDataModel{
		ID:                        plugin.TraceID,
		Type:                      traceType,
		AgentVersion:              plugin.AgentVersion,
		DataModelVersion:          plugin.DataModelVersion,
		ApplicationID:             plugin.ApplicationID,
		ApplicationDomainName:     plugin.ApplicationDomainName,
		ApplicationClassName:      plugin.ApplicationClassName,
		ApplicationName:           plugin.FunctionName,
		ApplicationVersion:        plugin.ApplicationVersion,
		ApplicationStage:          plugin.ApplicationStage,
		ApplicationRuntime:        plugin.ApplicationRuntime,
		ApplicationRuntimeVersion: plugin.ApplicationRuntimeVersion,
		ApplicationTags:           map[string]interface{}{},

		RootSpanID:      tr.rootSpan.Context().(ttracer.SpanContext).SpanID,
		StartTimestamp:  tr.data.startTime,
		FinishTimestamp: tr.data.finishTime,
		Duration:        tr.data.duration,
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
	tags[plugin.AwsLambdaInvocationTimeout] = tr.data.timeout

	// If this is the first invocation, it is a cold start
	if invocationCount == 1 {
		tags[plugin.AwsLambdaInvocationColdStart] = true
	} else {
		tags[plugin.AwsLambdaInvocationColdStart] = false
	}

	if tr.data.panicInfo != nil {
		tags[plugin.AwsError] = true
		tags[plugin.AwsErrorKind] = tr.data.panicInfo.Kind
		tags[plugin.AwsErrorMessage] = tr.data.panicInfo.Message
		tags[plugin.AwsErrorStack] = tr.data.panicInfo.Stack
	} else if tr.data.errorInfo != nil {
		tags[plugin.AwsError] = true
		tags[plugin.AwsErrorKind] = tr.data.errorInfo.Kind
		tags[plugin.AwsErrorMessage] = tr.data.errorInfo.Message
	}
	return tags
}

type spanDataModel struct {
	//Base fields
	ID                        string                 `json:"id"`
	Type                      string                 `json:"type"`
	AgentVersion              string                 `json:"agentVersion"`
	DataModelVersion          string                 `json:"dataModelVersion"`
	ApplicationID             string                 `json:"applicationId"`
	ApplicationDomainName     string                 `json:"applicationDomainName"`
	ApplicationClassName      string                 `json:"applicationClassName"`
	ApplicationName           string                 `json:"applicationName"`
	ApplicationVersion        string                 `json:"applicationVersion"`
	ApplicationStage          string                 `json:"applicationStage"`
	ApplicationRuntime        string                 `json:"applicationRuntime"`
	ApplicationRuntimeVersion string                 `json:"applicationRuntimeVersion"`
	ApplicationTags           map[string]interface{} `json:"applicationTags"`

	TraceID       string `json:"traceId"`
	TransactionID string `json:"transactionId"`
	ParentSpanID  string `json:"parentSpanId"`

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

func (tr *trace) prepareSpanDataModel(ctx context.Context, span *ttracer.RawSpan) spanDataModel {
	return spanDataModel{
		ID:                        span.Context.SpanID,
		Type:                      spanType,
		AgentVersion:              plugin.AgentVersion,
		DataModelVersion:          plugin.DataModelVersion,
		ApplicationID:             plugin.ApplicationID,
		ApplicationDomainName:     plugin.ApplicationDomainName,
		ApplicationClassName:      plugin.ApplicationClassName,
		ApplicationName:           plugin.FunctionName,
		ApplicationVersion:        plugin.ApplicationVersion,
		ApplicationStage:          plugin.ApplicationStage,
		ApplicationRuntime:        plugin.ApplicationRuntime,
		ApplicationRuntimeVersion: plugin.ApplicationRuntimeVersion,
		ApplicationTags:           map[string]interface{}{},

		TraceID:       span.Context.TraceID,
		TransactionID: plugin.TransactionID,
		ParentSpanID:  span.ParentSpanID,

		DomainName:    span.DomainName,
		ClassName:     span.ClassName,
		ServiceName:   plugin.FunctionName, //TODO implement it with Opentracing
		OperationName: span.OperationName,

		StartTimestamp:  span.StartTimestamp,
		FinishTimestamp: span.EndTimestamp,
		Duration:        span.Duration(),
		Tags:            span.Tags,
		Logs:            map[string]spanLog{}, // TO DO get logs
	}
}