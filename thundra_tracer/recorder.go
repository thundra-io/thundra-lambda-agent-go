package ttracer

import (
	"sync"
)

// ThundraRecorder handles all of the `RawSpan` data generated via an
// associated `Tracer`
type ThundraRecorder interface {
	Record(event SpanEvent, span *RawSpan)
	GetRootSpan() *RawSpan
	Reset()
}

// TreeSpanRecorder is a simple thread-safe implementation of
// SpanRecorder that stores all reported spans in a tree and pushes last started span into stack.
type recorderImpl struct {
	sync.RWMutex
	activeStack        *spanStack
	finishedStack	   *spanStack
}

// NewThundraRecorder creates new recorder to use in tracer
func NewThundraRecorder() *recorderImpl {
	r := new(recorderImpl)
	return r
}

// Record is called when a new span is started.
// When a span starts, a spantree is created and holds information about span and its children spans.
// Then it pushes this spantree into stack to actively hold information about which spantree is currently running.
func (r *recorderImpl) Record(event SpanEvent, span *RawSpan) {
	r.Lock()
	defer r.Unlock()

	if event == SpanStartEvent {
		r.recordStartSpan(span *RawSpan)
	} else if event == SpanFinishEvent {
		r.recordFinishSpan(span *RawSpan)
	}
}

func (r *recorderImpl) recordStartSpan(span *RawSpan) {
	r.activeStack.Push(span)
}

func (r *recorderImpl) recordFinishSpan(span *RawSpan) {
	r.activeStack.Pop()
	r.finishedStack.Push(span)
}

// Reset flushes data
func (r *recorderImpl) Reset() {
	r.Lock()
	defer r.Unlock()
	r.finishedStack.Clear()
}
