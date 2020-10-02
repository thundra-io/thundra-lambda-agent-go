package invocation

import (
	"context"

	"github.com/thundra-io/thundra-lambda-agent-go/v2/tracer"

	"github.com/thundra-io/thundra-lambda-agent-go/v2/application"
	"github.com/thundra-io/thundra-lambda-agent-go/v2/constants"
	"github.com/thundra-io/thundra-lambda-agent-go/v2/plugin"
	"github.com/thundra-io/thundra-lambda-agent-go/v2/utils"
)

// invocationPlugin is the simplest form of data collected from lambda functions. It is collected for any case.
type invocationDataModel struct {
	//Base fields
	plugin.BaseDataModel
	ID                  string                 `json:"id"`
	Type                string                 `json:"type"`
	TraceID             string                 `json:"traceId"`
	TransactionID       string                 `json:"transactionId"`
	SpanID              string                 `json:"spanId"`
	ApplicationPlatform string                 `json:"applicationPlatform"`
	FunctionRegion      string                 `json:"functionRegion"`
	StartTimestamp      int64                  `json:"startTimestamp"`  // Invocation start time in UNIX Epoch milliseconds
	FinishTimestamp     int64                  `json:"finishTimestamp"` // Invocation end time in UNIX Epoch milliseconds
	Duration            int64                  `json:"duration"`        // Invocation time in milliseconds
	Erroneous           bool                   `json:"erroneous"`       // Shows if the invocationPlugin failed with an error
	ErrorType           string                 `json:"errorType"`       // Type of the thrown error
	ErrorMessage        string                 `json:"errorMessage"`    // Message of the thrown error
	ErrorCode           string                 `json:"errorCode"`       // Numeric code of the error, such as 404 for HttpError
	ColdStart           bool                   `json:"coldStart"`       // Shows if the invocationPlugin is cold started
	Timeout             bool                   `json:"timeout"`         // Shows if the invocationPlugin is timed out
	Tags                map[string]interface{} `json:"tags"`
	UserTags            map[string]interface{} `json:"userTags"`
	IncomingTraceLinks  []string               `json:"incomingTraceLinks"`
	OutgoingTraceLinks  []string               `json:"outgoingTraceLinks"`
	Resources           []Resource             `json:"resources"`
}

func (ip *invocationPlugin) prepareData(ctx context.Context) invocationDataModel {
	spanID := ""
	if ip.rootSpan != nil {
		spanID = ip.rootSpan.Context().(tracer.SpanContext).SpanID
	}
	tags := ip.prepareTags(ctx)

	return invocationDataModel{
		BaseDataModel:       plugin.GetBaseData(),
		ID:                  utils.GenerateNewID(),
		Type:                invocationType,
		TraceID:             plugin.TraceID,
		TransactionID:       plugin.TransactionID,
		SpanID:              spanID,
		ApplicationPlatform: constants.ApplicationPlatform,
		FunctionRegion:      application.FunctionRegion,
		StartTimestamp:      ip.data.startTimestamp,
		FinishTimestamp:     ip.data.finishTimestamp,
		Duration:            ip.data.duration,
		Erroneous:           ip.data.erroneous,
		ErrorType:           ip.data.errorType,
		ErrorMessage:        ip.data.errorMessage,
		ErrorCode:           ip.data.errorCode,
		ColdStart:           ip.data.coldStart,
		Timeout:             ip.data.timeout,
		IncomingTraceLinks:  getIncomingTraceLinks(),
		OutgoingTraceLinks:  getOutgoingTraceLinks(),
		Tags:                tags,
		UserTags:            userInvocationTags,
		Resources:           getResources(spanID),
	}
}

func (ip *invocationPlugin) prepareTags(ctx context.Context) map[string]interface{} {
	tags := invocationTags

	// Put error related tags
	if ip.data.erroneous {
		tags["error"] = true
		tags["error.kind"] = ip.data.errorType
		tags["error.message"] = ip.data.errorMessage
	}
	arn := application.GetInvokedFunctionArn(ctx)
	xrayTraceID, xraySegmentID := utils.GetXRayTraceInfo(ctx)
	if len(xrayTraceID) > 0 {
		tags[constants.AwsXRayTraceID] = xrayTraceID
	}
	if len(xraySegmentID) > 0 {
		tags[constants.AwsXRaySegmentID] = xraySegmentID
	}
	tags[constants.AwsLambdaARN] = arn
	tags[constants.AwsAccountNo] = application.GetAwsAccountNo(arn)
	tags[constants.AwsLambdaInvocationColdStart] = ip.data.coldStart
	tags[constants.AwsLambdaInvocationRequestId] = application.GetAwsRequestID(ctx)
	tags[constants.AwsLambdaLogGroupName] = application.LogGroupName
	tags[constants.AwsLambdaLogStreamName] = application.LogStreamName
	tags[constants.AwsLambdaMemoryLimit] = application.MemoryLimit
	tags[constants.AwsLambdaMemoryUsage] = application.MemoryUsed
	tags[constants.AwsLambdaName] = application.ApplicationName
	tags[constants.AwsRegion] = application.FunctionRegion
	return tags
}
