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
	tags := tr.prepareTags(ctx, request, response)
	return traceData{
		Id:                        plugin.Id,
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

func (tr *trace) prepareTags(ctx context.Context, request json.RawMessage, response interface{}) map[string]interface{} {
	tags := map[string]interface{}{}
	tags[awsLambdaInvocationRequestId] = plugin.GetAwsRequestID(ctx)

	// If the agent's user doesn't want to send their request and response data, hide them.
	if !shouldHideRequest() {
		tags[awsLambdaInvocationRequest] = string(request)
	}
	if !shouldHideResponse() {
		tags[awsLambdaInvocationResponse] = response
	}

	tags[awsLambdaARN] = plugin.GetInvokedFunctionArn(ctx)
	tags[awsLambdaLogGroupName] = plugin.LogGroupName
	tags[awsLambdaLogStreamName] = plugin.LogStreamName
	tags[awsLambdaMemoryLimit] = plugin.MemoryLimit
	tags[awsLambdaName] = plugin.FunctionName
	tags[awsRegion] = plugin.FunctionRegion
	tags[awsLambdaInvocationTimeout] = tr.span.timeout

	// If this is the first invocation, it is a cold start
	if invocationCount == 1 {
		tags[awsLambdaInvocationColdStart] = true
	} else {
		tags[awsLambdaInvocationColdStart] = false
	}

	if tr.span.panicInfo != nil {
		tags[awsError] = true
		tags[awsErrorKind] = tr.span.panicInfo.Kind
		tags[awsErrorMessage] = tr.span.panicInfo.Message
		tags[awsErrorStack] = tr.span.panicInfo.Stack
	} else if tr.span.errorInfo != nil {
		tags[awsError] = true
		tags[awsErrorKind] = tr.span.errorInfo.Kind
		tags[awsErrorMessage] = tr.span.errorInfo.Message
	}

	return tags
}
