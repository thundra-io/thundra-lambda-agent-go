package invocation

import (
	"context"

	"github.com/thundra-io/thundra-lambda-agent-go/application"
	"github.com/thundra-io/thundra-lambda-agent-go/constants"
	"github.com/thundra-io/thundra-lambda-agent-go/plugin"
	"github.com/thundra-io/thundra-lambda-agent-go/utils"
)

// invocationPlugin is the simplest form of data collected from lambda functions. It is collected for any case.
type invocationDataModel struct {
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
	TraceID                   string                 `json:"traceId"`
	TransactionID             string                 `json:"transactionId"`
	SpanID                    string                 `json:"spanId"`
	FunctionPlatform          string                 `json:"functionPlatform"`
	FunctionName              string                 `json:"functionName"`
	FunctionRegion            string                 `json:"functionRegion"`
	StartTimestamp            int64                  `json:"startTimestamp"`  // Invocation start time in UNIX Epoch milliseconds
	FinishTimestamp           int64                  `json:"finishTimestamp"` // Invocation end time in UNIX Epoch milliseconds
	Duration                  int64                  `json:"duration"`        // Invocation time in milliseconds
	Erroneous                 bool                   `json:"erroneous"`       // Shows if the invocationPlugin failed with an error
	ErrorType                 string                 `json:"errorType"`       // Type of the thrown error
	ErrorMessage              string                 `json:"errorMessage"`    // Message of the thrown error
	ErrorCode                 string                 `json:"errorCode"`       // Numeric code of the error, such as 404 for HttpError
	ColdStart                 bool                   `json:"coldStart"`       // Shows if the invocationPlugin is cold started
	Timeout                   bool                   `json:"timeout"`         // Shows if the invocationPlugin is timed out
	Tags                      map[string]interface{} `json:"tags"`
}

func (ip *invocationPlugin) prepareData(ctx context.Context) invocationDataModel {
	tags := ip.prepareTags(ctx)
	return invocationDataModel{
		ID:                        utils.GenerateNewID(),
		Type:                      invocationType,
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
		TraceID:                   plugin.TraceID,
		TransactionID:             plugin.TransactionID,
		SpanID:                    "", // Optional Field
		FunctionPlatform:          constants.AwsFunctionPlatform,
		FunctionName:              application.ApplicationName,
		FunctionRegion:            application.FunctionRegion,
		StartTimestamp:            ip.data.startTimestamp,
		FinishTimestamp:           ip.data.finishTimestamp,
		Duration:                  ip.data.duration,
		Erroneous:                 ip.data.erroneous,
		ErrorType:                 ip.data.errorType,
		ErrorMessage:              ip.data.errorMessage,
		ErrorCode:                 ip.data.errorCode,
		ColdStart:                 ip.data.coldStart,
		Timeout:                   ip.data.timeout,
		Tags:                      tags,
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
	tags[constants.AwsLambdaARN] = application.GetInvokedFunctionArn(ctx)
	tags[constants.AwsLambdaInvocationColdStart] = ip.data.coldStart
	tags[constants.AwsLambdaInvocationRequestId] = application.GetAwsRequestID(ctx)
	tags[constants.AwsLambdaLogGroupName] = application.LogGroupName
	tags[constants.AwsLambdaLogStreamName] = application.LogStreamName
	tags[constants.AwsLambdaMemoryLimit] = application.MemoryLimit
	tags[constants.AwsLambdaName] = application.ApplicationName
	tags[constants.AwsRegion] = application.FunctionRegion
	return tags
}
