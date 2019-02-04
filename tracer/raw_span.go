package tracer

import (
	"strings"
	"time"

	"github.com/thundra-io/thundra-lambda-agent-go/ext"

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
	OperationName string

	StartTimestamp int64
	EndTimestamp   int64

	// operationGroup is lambda execution group
	DomainName string

	// operationType is lambda execution type
	ClassName string

	// Essentially an extension mechanism. Can be used for many purposes,
	// not to be enumerated here.
	Tags ot.Tags

	// The span's "microlog".
	Logs []ot.LogRecord
}

// Duration calculates the spans duration
func (s *RawSpan) Duration() int64 {
	if s.EndTimestamp != 0 {
		return s.EndTimestamp - s.StartTimestamp
	}

	return time.Now().Unix() - s.StartTimestamp
}

// GetTags filters the thundra tags and returns the remainings
func (s *RawSpan) GetTags() ot.Tags {
	ft := ot.Tags{}

	for k, v := range s.Tags {
		if !strings.HasPrefix(k, ext.ThundraTagPrefix) {
			ft[k] = v
		}
	}

	return ft
}
