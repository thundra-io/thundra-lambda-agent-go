package ttracer

import (
	"sync"
)

// SpanRecorder handles all of the `RawSpan` data generated via an
// associated `Tracer`
type SpanRecorder interface {
	GetSpans() []*RawSpan
	RecordSpan(span *RawSpan)
	Reset()
}

// InMemorySpanRecorder stores spans using a slice in a thread-safe way
type InMemorySpanRecorder struct {
	sync.RWMutex
	spans []*RawSpan
}

// NewInMemoryRecorder creates new InMemorySpanRecorder
func NewInMemoryRecorder() *InMemorySpanRecorder {
	return new(InMemorySpanRecorder)
}

// RecordSpan implements the respective method of SpanRecorder.
func (r *InMemorySpanRecorder) RecordSpan(span *RawSpan) {
	r.Lock()
	defer r.Unlock()
	r.spans = append(r.spans, span)
}

// GetSpans returns a copy of the array of spans accumulated so far.
func (r *InMemorySpanRecorder) GetSpans() []*RawSpan {
	r.RLock()
	defer r.RUnlock()
	spans := make([]*RawSpan, len(r.spans))
	copy(spans, r.spans)
	return spans
}

// Reset clears the internal array of spans.
func (r *InMemorySpanRecorder) Reset() {
	r.Lock()
	defer r.Unlock()
	r.spans = nil
}
