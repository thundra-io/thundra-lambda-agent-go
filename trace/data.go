package trace

import (
	"context"
	"encoding/json"

	"github.com/thundra-io/thundra-lambda-agent-go/config"

	"github.com/thundra-io/thundra-lambda-agent-go/application"
	"github.com/thundra-io/thundra-lambda-agent-go/constants"
	"github.com/thundra-io/thundra-lambda-agent-go/tracer"

	"github.com/thundra-io/thundra-lambda-agent-go/plugin"
)

type traceDataModel struct {
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
	RootSpanID                string                 `json:"rootSpanId"`
	StartTimestamp            int64                  `json:"startTimestamp"`
	FinishTimestamp           int64                  `json:"finishTimestamp"`
	Duration                  int64                  `json:"duration"`
	Tags                      map[string]interface{} `json:"tags"`
}

func (tr *tracePlugin) prepareTraceDataModel(ctx context.Context, request json.RawMessage, response interface{}) traceDataModel {
	tags := tr.prepareTraceTags(ctx, request, response)
	return traceDataModel{
		ID:                        plugin.TraceID,
		Type:                      traceType,
		AgentVersion:              constants.AgentVersion,
		DataModelVersion:          constants.DataModelVersion,
		ApplicationID:             application.ApplicationID,
		ApplicationDomainName:     application.ApplicationDomainName,
		ApplicationClassName:      application.ApplicationClassName,
		ApplicationName:           application.ApplicationName,
		ApplicationVersion:        application.ApplicationVersion,
		ApplicationStage:          application.ApplicationStage,
		ApplicationRuntime:        application.ApplicationRuntime,
		ApplicationRuntimeVersion: application.ApplicationRuntimeVersion,
		ApplicationTags:           application.ApplicationTags,
		RootSpanID:                tr.rootSpan.Context().(tracer.SpanContext).SpanID,
		StartTimestamp:            tr.data.startTime,
		FinishTimestamp:           tr.data.finishTime,
		Duration:                  tr.data.duration,
		Tags:                      tags,
	}
}

func (tr *tracePlugin) prepareTraceTags(ctx context.Context, request json.RawMessage, response interface{}) map[string]interface{} {
	tags := map[string]interface{}{}
	tags[constants.AwsLambdaInvocationRequestId] = application.GetAwsRequestID(ctx)

	// If the agent's user doesn't want to send their request and response data, hide them.
	if !config.TraceRequestDisabled {
		tags[constants.AwsLambdaInvocationRequest] = string(request)
	}
	if !config.TraceResponseDisabled {
		tags[constants.AwsLambdaInvocationResponse] = response
	}

	tags[constants.AwsLambdaARN] = application.GetInvokedFunctionArn(ctx)
	tags[constants.AwsLambdaLogGroupName] = application.LogGroupName
	tags[constants.AwsLambdaLogStreamName] = application.LogStreamName
	tags[constants.AwsLambdaMemoryLimit] = application.MemoryLimit
	tags[constants.AwsLambdaName] = application.ApplicationName
	tags[constants.AwsRegion] = application.FunctionRegion
	tags[constants.AwsLambdaInvocationTimeout] = tr.data.timeout

	// If this is the first invocation, it is a cold start
	if invocationCount == 1 {
		tags[constants.AwsLambdaInvocationColdStart] = true
	} else {
		tags[constants.AwsLambdaInvocationColdStart] = false
	}

	if tr.data.panicInfo != nil {
		tags[constants.AwsError] = true
		tags[constants.AwsErrorKind] = tr.data.panicInfo.Kind
		tags[constants.AwsErrorMessage] = tr.data.panicInfo.Message
		tags[constants.AwsErrorStack] = tr.data.panicInfo.Stack
	} else if tr.data.errorInfo != nil {
		tags[constants.AwsError] = true
		tags[constants.AwsErrorKind] = tr.data.errorInfo.Kind
		tags[constants.AwsErrorMessage] = tr.data.errorInfo.Message
	}
	return tags
}

type spanDataModel struct {
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
	TraceID                   string                 `json:"traceId"`
	TransactionID             string                 `json:"transactionId"`
	ParentSpanID              string                 `json:"parentSpanId"`
	SpanOrder                 int64                  `json:"spanOrder"`
	DomainName                string                 `json:"domainName"`
	ClassName                 string                 `json:"className"`
	ServiceName               string                 `json:"serviceName"`
	OperationName             string                 `json:"operationName"`
	StartTimestamp            int64                  `json:"startTimestamp"`
	FinishTimestamp           int64                  `json:"finishTimestamp"`
	Duration                  int64                  `json:"duration"`
	Tags                      map[string]interface{} `json:"tags"`
	Logs                      map[string]spanLog     `json:"logs"`
}

type spanLog struct {
	Name      string      `json:"name"`
	Value     interface{} `json:"value"`
	Timestamp int64       `json:"timestamp"`
}

func (tr *tracePlugin) prepareSpanDataModel(ctx context.Context, span *tracer.RawSpan) spanDataModel {
	return spanDataModel{
		ID:                        span.Context.SpanID,
		Type:                      spanType,
		AgentVersion:              constants.AgentVersion,
		DataModelVersion:          constants.DataModelVersion,
		ApplicationID:             application.ApplicationID,
		ApplicationDomainName:     application.ApplicationDomainName,
		ApplicationClassName:      application.ApplicationClassName,
		ApplicationName:           application.ApplicationName,
		ApplicationVersion:        application.ApplicationVersion,
		ApplicationStage:          application.ApplicationStage,
		ApplicationRuntime:        application.ApplicationRuntime,
		ApplicationRuntimeVersion: application.ApplicationRuntimeVersion,
		ApplicationTags:           application.ApplicationTags,
		TraceID:                   span.Context.TraceID,
		TransactionID:             span.Context.TransactionID,
		ParentSpanID:              span.ParentSpanID,
		DomainName:                span.DomainName,
		ClassName:                 span.ClassName,
		ServiceName:               application.ApplicationName,
		OperationName:             span.OperationName,
		StartTimestamp:            span.StartTimestamp,
		FinishTimestamp:           span.EndTimestamp,
		Duration:                  span.Duration(),
		Tags:                      span.GetTags(),
		Logs:                      map[string]spanLog{}, // TO DO get logs
	}
}
