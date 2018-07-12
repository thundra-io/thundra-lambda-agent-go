package thundra_tracer

import (
	"sync"
)

// A SpanRecorder handles all of the `RawSpan` data generated via an
// associated `Tracer`
type SpanRecorder interface {
	RecordSpanStarted(span *RawSpan)
	RecordSpanEnded()
	GetSpanTree() *RawSpanTree
	Reset()
}

// TreeSpanRecorder is a simple thread-safe implementation of
// SpanRecorder that stores all reported spans in a tree and pushes last started span into stack.
type TreeSpanRecorder struct {
	sync.RWMutex

	spanTree        *RawSpanTree
	activeSpanStack spanTreeStack
}

// NewTreeSpanRecorder creates new TreeSpanRecorder
func NewTreeSpanRecorder() *TreeSpanRecorder {
	r := new(TreeSpanRecorder)
	return r
}

// RecordSpanStarted is called when a new span is started.
// When a span starts, a spantree is created and holds information about span and its children spans.
// Then it pushes this spantree into stack to actively hold information about which spantree is currently running.
func (r *TreeSpanRecorder) RecordSpanStarted(span *RawSpan) {
	r.Lock()
	defer r.Unlock()

	t := newRawSpanTree(span)
	if r.spanTree == nil {
		r.spanTree = t
		r.activeSpanStack.Push(t)
		return
	}
	top, err := r.activeSpanStack.Top()
	if err != nil {

	}

	top.addChild(t)
	r.activeSpanStack.Push(t)
}

// RecordSpanEnded is called when the span is finished.
// When a span finishes it is popped from the stack.
func (r *TreeSpanRecorder) RecordSpanEnded() {
	r.Lock()
	defer r.Unlock()

	r.activeSpanStack.Pop()
}

// GetSpanTree returns spanTree
func (r *TreeSpanRecorder) GetSpanTree() *RawSpanTree {
	r.Lock()
	defer r.Unlock()

	return r.spanTree
}

// Reset flushes data
func (r *TreeSpanRecorder) Reset() {
	r.Lock()
	defer r.Unlock()
	r.spanTree = nil
	r.activeSpanStack = *NewStack()
}
