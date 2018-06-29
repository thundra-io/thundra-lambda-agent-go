package mytracer

import (
	ot "github.com/opentracing/opentracing-go"
	"time"
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

func (t *tracerImpl) StartSpan(operationName string, opts ...ot.StartSpanOption) ot.Span {
	sso := startSpanOptions{}
	for _, o := range opts {
		o.Apply(&sso.Options)
	}
	// Start time.
	startTime := sso.Options.StartTime
	if startTime.IsZero() {
		startTime = time.Now()
	}
	// Build the new span. This is the only allocation: We'll return this as
	// an opentracing.Span.
	sp := &spanImpl{}
	// Look for a parent in the list of References.
	//
	// TODO: would be nice if basictracer did something with all
	// References, not just the first one.
	ReferencesLoop:
		for _, ref := range opts.References {
			switch ref.Type {
			case opentracing.ChildOfRef,
				opentracing.FollowsFromRef:

				refCtx := ref.ReferencedContext.(SpanContext)
				sp.raw.Context.TraceID = refCtx.TraceID
				sp.raw.Context.SpanID = randomID()
				sp.raw.Context.Sampled = refCtx.Sampled
				sp.raw.ParentSpanID = refCtx.SpanID

				if l := len(refCtx.Baggage); l > 0 {
					sp.raw.Context.Baggage = make(map[string]string, l)
					for k, v := range refCtx.Baggage {
						sp.raw.Context.Baggage[k] = v
					}
				}
				break ReferencesLoop
			}
		}
	/*if sp.raw.Context.TraceID == 0 {
		// No parent Span found; allocate new trace and span ids and determine
		// the Sampled status.
		// TODO change it for generation
		sp.raw.Context.TraceID, sp.raw.Context.SpanID = 1, 2
	}*/
	sp.tracer = t
	sp.raw.Operation = operationName
	sp.raw.Start = startTime
	sp.raw.Duration = -1
	//sp.raw.Tags = opts.Tags
	return sp
}

//
func (tracer *tracerImpl) Inject(sc ot.SpanContext, format interface{}, carrier interface{}) error {
	return errors.New("Inject has not been supported yet")
}

func (tracer *tracerImpl) Extract(format interface{}, carrier interface{}) (ot.SpanContext, error) {
	return nil, errors.New("Extract has not been supported yet")
}
