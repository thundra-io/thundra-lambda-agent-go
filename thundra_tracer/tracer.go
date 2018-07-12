package thundra_tracer

import (
	ot "github.com/opentracing/opentracing-go"
	"github.com/pkg/errors"
)

// NewWithOptions creates a customized Tracer.
func NewWithOptions(opts Options) ot.Tracer {
	rval := &tracerImpl{opts: opts}
	return rval
}

// New creates and returns a standard Tracer which defers completed Spans to
// `recorder`.
// Spans created by this Tracer support the ext.SamplingPriority tag: Setting
// ext.SamplingPriority causes the Span to be Sampled from that point on.
func New(recorder SpanRecorder) ot.Tracer {
	opts := DefaultOptions()
	opts.Recorder = recorder
	return NewWithOptions(opts)
}

type tracerImpl struct {
	opts Options
}

func (t *tracerImpl) StartSpan(operationName string, sso ...ot.StartSpanOption) ot.Span {
	return newSpan(operationName, t, sso)
}

// TODO Will be implemented
func (tracer *tracerImpl) Inject(sc ot.SpanContext, format interface{}, carrier interface{}) error {
	return errors.New("Inject has not been supported yet")
}

// TODO Will be implemented
func (tracer *tracerImpl) Extract(format interface{}, carrier interface{}) (ot.SpanContext, error) {
	return nil, errors.New("Extract has not been supported yet")
}
