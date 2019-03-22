package trace

import (
	"context"
	"encoding/json"

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
		RootSpanID:                tr.RootSpan.Context().(tracer.SpanContext).SpanID,
		StartTimestamp:            tr.Data.StartTime,
		FinishTimestamp:           tr.Data.FinishTime,
		Duration:                  tr.Data.Duration,
	}
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
	// If a span have no rootSpanID (other than the root span)
	// Set rootSpan's ID as the parent ID for that span
	rootSpanID := tr.RootSpan.Context().(tracer.SpanContext).SpanID
	if len(span.ParentSpanID) == 0 && span.Context.SpanID != rootSpanID {
		span.ParentSpanID = rootSpanID
	}
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
