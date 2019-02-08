package tracer

import (
	"strings"
	"time"

	"github.com/thundra-io/thundra-lambda-agent-go/ext"

	ot "github.com/opentracing/opentracing-go"
)

// RawSpan encapsulates all state associated with a (finished) Span.
type RawSpan struct {
	Context        SpanContext
	ParentSpanID   string
	OperationName  string
	StartTimestamp int64
	EndTimestamp   int64
	DomainName     string
	ClassName      string
	Tags           ot.Tags
	Logs           []ot.LogRecord
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
