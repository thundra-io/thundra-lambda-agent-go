package trace

import (
	"context"
	"encoding/json"

	"github.com/thundra-io/thundra-lambda-agent-go/v2/application"
	"github.com/thundra-io/thundra-lambda-agent-go/v2/tracer"

	"github.com/thundra-io/thundra-lambda-agent-go/v2/plugin"
)

type traceDataModel struct {
	plugin.BaseDataModel
	ID              string                 `json:"id"`
	Type            string                 `json:"type"`
	RootSpanID      string                 `json:"rootSpanId"`
	StartTimestamp  int64                  `json:"startTimestamp"`
	FinishTimestamp int64                  `json:"finishTimestamp"`
	Duration        int64                  `json:"duration"`
	Tags            map[string]interface{} `json:"tags"`
}

func (tr *tracePlugin) prepareTraceDataModel(ctx context.Context, request json.RawMessage, response interface{}) traceDataModel {
	return traceDataModel{
		BaseDataModel:   plugin.GetBaseData(),
		ID:              plugin.TraceID,
		Type:            traceType,
		RootSpanID:      tr.RootSpan.Context().(tracer.SpanContext).SpanID,
		StartTimestamp:  tr.Data.StartTime,
		FinishTimestamp: tr.Data.FinishTime,
		Duration:        tr.Data.Duration,
	}
}

type spanDataModel struct {
	plugin.BaseDataModel
	ID              string                 `json:"id"`
	Type            string                 `json:"type"`
	TraceID         string                 `json:"traceId"`
	TransactionID   string                 `json:"transactionId"`
	ParentSpanID    string                 `json:"parentSpanId"`
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

func (tr *tracePlugin) prepareSpanDataModel(ctx context.Context, span *tracer.RawSpan) spanDataModel {
	// If a span have no rootSpanID (other than the root span)
	// Set rootSpan's ID as the parent ID for that span
	rootSpanID := tr.RootSpan.Context().(tracer.SpanContext).SpanID
	if len(span.ParentSpanID) == 0 && span.Context.SpanID != rootSpanID {
		span.ParentSpanID = rootSpanID
	}
	return spanDataModel{
		BaseDataModel:   plugin.GetBaseData(),
		ID:              span.Context.SpanID,
		Type:            spanType,
		TraceID:         span.Context.TraceID,
		TransactionID:   span.Context.TransactionID,
		ParentSpanID:    span.ParentSpanID,
		DomainName:      span.DomainName,
		ClassName:       span.ClassName,
		ServiceName:     application.ApplicationName,
		OperationName:   span.OperationName,
		StartTimestamp:  span.StartTimestamp,
		FinishTimestamp: span.EndTimestamp,
		Duration:        span.Duration(),
		Tags:            span.GetTags(),
		Logs:            map[string]spanLog{}, // TO DO get logs
	}
}
