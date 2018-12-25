package thundra_tracer

import (
	ot "github.com/opentracing/opentracing-go"
)

// RawSpan encapsulates all state associated with a (finished) Span.
type RawSpan struct {
	// Those recording the RawSpan should also record the contents of its
	// SpanContext.
	Context SpanContext

	// The SpanID of this SpanContext's first intra-trace reference (i.e.,
	// "parent"), or 0 if there is no parent.
	ParentSpanID string

	// The name of the "operation" this span is an instance of. (Called a "span
	// name" in some implementations)
	Operation string

	TraceId 	   string
	SpanOrder	   int64
	DomainName 	   string
	ClassName      string
	ServiceName    string
	OperationName  string
	StartTimestamp int64
	EndTimestamp   int64
	Duration       int64

	// operationGroup is lambda execution group
	operationGroup operationGroup

	// operationType is lambda execution type
	operationType operationType

	// Essentially an extension mechanism. Can be used for many purposes,
	// not to be enumerated here.
	Tags ot.Tags

	// The span's "microlog".
	Logs []ot.LogRecord
}

type SpanData struct {
	TraceId          string                 		`json:"traceId"`
	TransactionId    string							`json:"transactionId"`
	SpanOrder		 int64							`json:"spanOrder"`
	DomainName		 string							`json:"domainName"`
	ClassName		 string							`json:"className"`
	ServiceName		 string							`json:"serviceName"`
	OperationName	 string							`json:"operationName"`
	StartTimeStamp   int64							`json:"startTimeStamp"`
	FinishTimeStamp  int64							`json:"finishTimeStamp"`
	Tags 			 map[string]interface{} 		`json:"tags"`
	Logs 			 map[string]interface{} 		`json:"logs"`

}

func PrepareSpanData(raw RawSpan) SpanData{
	return SpanData{
		TraceId: raw.TraceId,
		TransactionId: "",
	}
}