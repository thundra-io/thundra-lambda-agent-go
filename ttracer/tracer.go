package ttracer

import (
	ot "github.com/opentracing/opentracing-go"
	"github.com/pkg/errors"
	uuid "github.com/satori/go.uuid"
	"github.com/thundra-io/thundra-lambda-agent-go/plugin"
)

// NewWithOptions creates a customized Tracer.
func NewWithOptions(opts Options) ot.Tracer {
	rval := &tracerImpl{options: opts}
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
	options Options
}

// StartSpan starts a new span with options and returns it.
func (t *tracerImpl) StartSpan(operationName string, opts ...ot.StartSpanOption) ot.Span {
	sso := ot.StartSpanOptions{}
	for _, o := range opts {
		o.Apply(&sso)
	}

	return t.StartSpanWithOptions(operationName, sso)
}

func (t *tracerImpl) StartSpanWithOptions(operationName string, opts ot.StartSpanOptions) ot.Span {
	tags := opts.Tags
	if tags == nil {
		tags = map[string]interface{}{}
	}

	newSpan := t.getSpan()
	for _, ref := range opts.References {
		if ref.Type == ot.ChildOfRef {
			refSpanCtx := ref.ReferencedContext.(SpanContext)
			newSpan.raw.Context.TraceID = refSpanCtx.TraceID
			newSpan.raw.Context.SpanID = uuid.NewV4().String()
			newSpan.raw.ParentSpanID = refSpanCtx.SpanID

			if l := len(refSpanCtx.Baggage); l > 0 {
				newSpan.raw.Context.Baggage = make(map[string]string, l)
				for k, v := range refSpanCtx.Baggage {
					newSpan.raw.Context.Baggage[k] = v
				}
			}
		}
	}

	if newSpan.raw.Context.TraceID == "" {
		// Couldn't find a parent span then create a new trace and span id
		newSpan.raw.Context.TraceID = uuid.NewV4().String()
		newSpan.raw.Context.SpanID = uuid.NewV4().String()
	}

	newSpan.tracer = t
	newSpan.raw.OperationName = operationName
	newSpan.raw.ClassName = plugin.DefaultClassName
	newSpan.raw.DomainName = plugin.DefaultDomainName
	newSpan.raw.StartTimestamp = GetTimestamp()
	newSpan.raw.Tags = tags
	newSpan.raw.Logs = []ot.LogRecord{}

	// Add to recorder
	t.options.Recorder.RecordSpan(&newSpan.raw)
	return newSpan
}

func (t *tracerImpl) getSpan() *spanImpl {
	return &spanImpl{}
}

// TODO Will be implemented
func (t *tracerImpl) Inject(sc ot.SpanContext, format interface{}, carrier interface{}) error {
	return errors.New("Inject has not been supported yet")
}

// TODO Will be implemented
func (t *tracerImpl) Extract(format interface{}, carrier interface{}) (ot.SpanContext, error) {
	return nil, errors.New("Extract has not been supported yet")
}
