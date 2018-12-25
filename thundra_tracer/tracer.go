package ttracer

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
func New(recorder SpanRecorder) ot.Tracer {
	opts := DefaultOptions()
	opts.Recorder = recorder
	return NewWithOptions(opts)
}

type tracerImpl struct {
	opts Options
}

// StartSpan starts a new span with options and returns it.
func (t *tracerImpl) StartSpan(operationName string, sso ...ot.StartSpanOption) ot.Span {
	return newSpan(operationName, t, sso)
}

// TODO Will be implemented
func (tracer *tracerImpl) Inject(sc ot.SpanContext, format interface{}, carrier interface{}) error {
	panic(errors.New("Inject has not been supported yet"))
	return errors.New("Inject has not been supported yet")
}

// TODO Will be implemented
func (tracer *tracerImpl) Extract(format interface{}, carrier interface{}) (ot.SpanContext, error) {
	panic(errors.New("Inject has not been supported yet"))
	return nil, errors.New("Extract has not been supported yet")
}
