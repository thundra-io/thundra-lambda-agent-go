package tracer

import (
	ot "github.com/opentracing/opentracing-go"
	"github.com/pkg/errors"
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
	sso := ot.StartSpanOptions{}
	for _, o := range opts {
		o.Apply(&sso)
	}
	return t.StartSpanWithOptions(operationName, sso)
}

func (t *tracerImpl) StartSpanWithOptions(operationName string, opts ot.StartSpanOptions) ot.Span {
	newSpan := t.getSpan()

	newSpan.tracer = t
	newSpan.raw.Tags = opts.Tags
	newSpan.raw.Logs = []ot.LogRecord{}
	newSpan.raw.OperationName = operationName
	newSpan.raw.Context.TransactionID = plugin.TransactionID
	newSpan.raw.Context.TraceID = plugin.TraceID
	newSpan.raw.Context.SpanID = utils.GenerateNewID()

	for _, ref := range opts.References {
		if ref.Type == ot.ChildOfRef {
			parentCtx := ref.ReferencedContext.(SpanContext)
			newSpan.setParent(parentCtx)
		}
	}

	if opts.StartTime.IsZero() {
		newSpan.raw.StartTimestamp = utils.GetTimestamp()
	} else {
		newSpan.raw.StartTimestamp = utils.TimeToMs(opts.StartTime)
	}

	className, ok := opts.Tags[ext.ClassNameKey]
	if ok {
		classNameStr, ok := className.(string)
		if ok {
			newSpan.raw.ClassName = classNameStr
		} else {
			newSpan.raw.ClassName = constants.DefaultClassName
		}
	} else {
		newSpan.raw.ClassName = constants.DefaultClassName
	}

	domainName, ok := opts.Tags[ext.DomainNameKey]
	if ok {
		domainNameStr, ok := domainName.(string)
		if ok {
			newSpan.raw.DomainName = domainNameStr
		} else {
			newSpan.raw.DomainName = constants.DefaultDomainName
		}
	} else {
		newSpan.raw.DomainName = constants.DefaultDomainName
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
