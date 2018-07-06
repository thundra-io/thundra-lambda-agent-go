package otTracer

import (
	"sync"
)

// A SpanRecorder handles all of the `RawSpan` data generated via an
// associated `Tracer` (see `NewStandardTracer`) instance. It also names
// the containing process and provides access to a straightforward tag map.
type SpanRecorder interface {
	// Implementations must determine whether and where to store `span`.
	//RecordSpan(span RawSpan)

	RecordSpanStarted(span *RawSpan)
	RecordSpanEnded()
}

// InMemorySpanRecorder is a simple thread-safe implementation of
// SpanRecorder that stores all reported spans in memory, accessible
// via reporter.GetSpans(). It is primarily intended for testing purposes.
type InMemorySpanRecorder struct {
	sync.RWMutex

	allSpansTree     *RawSpanTree
	activeSpansStack spanTreeStack
}

// NewInMemoryRecorder creates new InMemorySpanRecorder
func NewInMemoryRecorder() *InMemorySpanRecorder {
	r := new(InMemorySpanRecorder)
	return r
}

/*// RecordSpan implements the respective method of SpanRecorder.
func (r *InMemorySpanRecorder) RecordSpan(span RawSpan) {
	r.Lock()
	defer r.Unlock()
	r.spans = append(r.spans, span)
}

// GetSpans returns a copy of the array of spans accumulated so far.
func (r *InMemorySpanRecorder) GetSpans() []RawSpan {
	r.RLock()
	defer r.RUnlock()
	spans := make([]RawSpan, len(r.spans))
	copy(spans, r.spans)
	return spans
}

// Reset clears the internal array of spans.
func (r *InMemorySpanRecorder) Reset() {
	r.Lock()
	defer r.Unlock()
	r.spans = nil
}*/

func (r *InMemorySpanRecorder) RecordSpanStarted(span *RawSpan) {
	r.Lock()
	defer r.Unlock()

	t := newRawSpanTree(span)
	if r.allSpansTree == nil {
		r.allSpansTree = t
		r.activeSpansStack.Push(t)
		return
	}
	top, err := r.activeSpansStack.Top()
	if err != nil {

	}

	top.addChild(t)
	r.activeSpansStack.Push(t)
}

func (r *InMemorySpanRecorder) RecordSpanEnded() {
	r.Lock()
	defer r.Unlock()

	r.activeSpansStack.Pop()
}

func (r *InMemorySpanRecorder) GetAllSpansTree() *RawSpanTree {
	r.Lock()
	defer r.Unlock()

	return r.allSpansTree
}

func (r *InMemorySpanRecorder) Reset() {
	r.Lock()
	defer r.Unlock()
	r.allSpansTree = nil
	r.activeSpansStack = *NewStack()
}
