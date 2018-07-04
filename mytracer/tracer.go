package mytracer

import (
	ot "github.com/opentracing/opentracing-go"
	"github.com/pkg/errors"
)

type tracer interface {
	ot.Tracer
}

// Options allows creating a customized Tracer via NewWithOptions.
type Options struct {
	// MaxLogsPerSpan limits the number of Logs in a span (if set to a nonzero
	// value). If a span has more logs than this value, logs are dropped as
	// necessary (and replaced with a log describing how many were dropped).
	//
	// About half of the MaxLogPerSpan logs kept are the oldest logs, and about
	// half are the newest logs.
	//
	// If NewSpanEventListener is set, the callbacks will still fire for all log
	// events. This value is ignored if DropAllLogs is true.
	MaxLogsPerSpan int
	// Recorder receives Spans which have been finished.
	Recorder SpanRecorder
	// DropAllLogs turns log events on all Spans into no-ops.
	// If NewSpanEventListener is set, the callbacks will still fire.
	DropAllLogs bool
}

// DefaultOptions returns an Options object with a 1 in 64 sampling rate and
// all options disabled. A Recorder needs to be set manually before using the
// returned object with a Tracer.
func DefaultOptions() Options {
	return Options{
		MaxLogsPerSpan: 100,
	}
}

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

//
func (tracer *tracerImpl) Inject(sc ot.SpanContext, format interface{}, carrier interface{}) error {
	return errors.New("Inject has not been supported yet")
}

func (tracer *tracerImpl) Extract(format interface{}, carrier interface{}) (ot.SpanContext, error) {
	return nil, errors.New("Extract has not been supported yet")
}
