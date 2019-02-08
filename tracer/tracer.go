package tracer

import (
	ot "github.com/opentracing/opentracing-go"
	"github.com/pkg/errors"
	uuid "github.com/satori/go.uuid"
	"github.com/thundra-io/thundra-lambda-agent-go/constants"
	"github.com/thundra-io/thundra-lambda-agent-go/ext"
	"github.com/thundra-io/thundra-lambda-agent-go/plugin"
	"github.com/thundra-io/thundra-lambda-agent-go/utils"
)

// New creates and returns a standard Tracer which defers completed Spans to
// `recorder`.
func New(recorder SpanRecorder) ot.Tracer {
	return &tracerImpl{
		Recorder: recorder,
	}
}

type tracerImpl struct {
	Recorder SpanRecorder
}

// StartSpan starts a new span with options and returns it.
func (t *tracerImpl) StartSpan(operationName string, opts ...ot.StartSpanOption) ot.Span {
	sso := ot.StartSpanOptions{
		Tags: make(map[string]interface{}),
	}
	for _, o := range opts {
		o.Apply(&sso)
	}
	return t.StartSpanWithOptions(operationName, sso)
}

func (t *tracerImpl) StartSpanWithOptions(operationName string, opts ot.StartSpanOptions) ot.Span {
	tags := opts.Tags

	newSpan := t.getSpan()
	newSpan.raw.Context.TraceID = plugin.TraceID
	newSpan.raw.Context.SpanID = uuid.NewV4().String()

	for _, ref := range opts.References {
		if ref.Type == ot.ChildOfRef {
			parentCtx := ref.ReferencedContext.(SpanContext)
			newSpan.setParent(parentCtx)
		}
	}

	newSpan.tracer = t
	newSpan.raw.OperationName = operationName
	newSpan.raw.StartTimestamp = utils.GetTimestamp()
	newSpan.raw.Tags = tags
	newSpan.raw.Logs = []ot.LogRecord{}

	className, ok := tags[ext.ClassNameKey]
	if !ok {
		newSpan.raw.ClassName = constants.DefaultClassName
	} else {
		newSpan.raw.ClassName = className.(string)
	}

	domainName, ok := tags[ext.DomainNameKey]
	if !ok {
		newSpan.raw.DomainName = constants.DefaultDomainName
	} else {
		newSpan.raw.DomainName = domainName.(string)
	}

	// Add to recorder
	t.Recorder.RecordSpan(&newSpan.raw)
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
