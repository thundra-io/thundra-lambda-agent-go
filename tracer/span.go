package tracer

import (
	"sync"
	"time"

	ot "github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/log"
	"github.com/thundra-io/thundra-lambda-agent-go/utils"
)

type spanImpl struct {
	tracer     *tracerImpl
	sync.Mutex // protects the fields below
	raw        RawSpan
	// The number of logs dropped because of MaxLogsPerSpan.
	numDroppedLogs int
}

// Finish is the last call that is made
func (s *spanImpl) Finish() {
	s.FinishWithOptions(ot.FinishOptions{})
}

// FinishWithOptions finishes span and adds the given options to it
func (s *spanImpl) FinishWithOptions(opts ot.FinishOptions) {
	if opts.FinishTime.IsZero() {
		s.raw.EndTimestamp = utils.GetTimestamp()
	} else {
		s.raw.EndTimestamp = utils.TimeToMs(opts.FinishTime)
	}

	s.Lock()
	defer s.Unlock()

	for _, lr := range opts.LogRecords {
		s.appendLog(lr)
	}
	for _, ld := range opts.BulkLogData {
		s.appendLog(ld.ToLogRecord())
	}
}

// Deprecated: use LogFields or LogKV.
func (s *spanImpl) LogEvent(event string) {
	s.Log(ot.LogData{
		Event: event,
	})
}

// Deprecated: use LogFields or LogKV.
func (s *spanImpl) LogEventWithPayload(event string, payload interface{}) {
	s.Log(ot.LogData{
		Event:   event,
		Payload: payload,
	})
}

// Log .
func (s *spanImpl) Log(ld ot.LogData) {
	s.Lock()
	defer s.Unlock()

	if ld.Timestamp.IsZero() {
		ld.Timestamp = time.Now()
	}

	s.appendLog(ld.ToLogRecord())
}

// SetOperationName sets operation name.
func (s *spanImpl) SetOperationName(operationName string) ot.Span {
	s.Lock()
	defer s.Unlock()
	s.raw.OperationName = operationName
	return s
}

// SetTag sets a tag with given key and value.
func (s *spanImpl) SetTag(key string, value interface{}) ot.Span {
	s.Lock()
	defer s.Unlock()
	if s.raw.Tags == nil {
		s.raw.Tags = ot.Tags{}
	}
	s.raw.Tags[key] = value
	return s
}

// LogKV logs logFields
func (s *spanImpl) LogKV(keyValues ...interface{}) {
	fields, err := log.InterleavedKVToFields(keyValues...)
	if err != nil {
		s.LogFields(log.Error(err), log.String("function", "LogKV"))
		return
	}
	s.LogFields(fields...)
}

func (s *spanImpl) appendLog(lr ot.LogRecord) {
	s.raw.Logs = append(s.raw.Logs, lr)
	return
}

// LogFields parses parameter fields sequentially, as first one is the key and the second is it's value.
func (s *spanImpl) LogFields(fields ...log.Field) {
	lr := ot.LogRecord{
		Fields: fields,
	}
	s.Lock()
	defer s.Unlock()
	if lr.Timestamp.IsZero() {
		lr.Timestamp = time.Now()
	}
	s.appendLog(lr)
}

// Tracer return span's Tracer
func (s *spanImpl) Tracer() ot.Tracer {
	return s.tracer
}

// Context returns SpanContext
func (s *spanImpl) Context() ot.SpanContext {
	return s.raw.Context
}

// SetBaggageItem sets BaggageItem
func (s *spanImpl) SetBaggageItem(key, val string) ot.Span {
	s.Lock()
	defer s.Unlock()
	s.raw.Context = s.raw.Context.WithBaggageItem(key, val)
	return s
}

// BaggageItem returns BaggageItem
func (s *spanImpl) BaggageItem(key string) string {
	s.Lock()
	defer s.Unlock()
	return s.raw.Context.Baggage[key]
}

// Operation returns the name of the "operation" this span is an instance of
func (s *spanImpl) OperationName() string {
	return s.raw.OperationName
}

// StartTimestamp returns StartTimestamp
func (s *spanImpl) StartTimestamp() int64 {
	return s.raw.StartTimestamp
}

// GetRaw casts opentracing span interface to spanImpl struct
func GetRaw(ots ot.Span) (*RawSpan, bool) {
	s, ok := ots.(*spanImpl)

	return &s.raw, ok
}

func (s *spanImpl) setParent(parentCtx SpanContext) {
	s.raw.ParentSpanID = parentCtx.SpanID

	if l := len(parentCtx.Baggage); l > 0 {
		s.raw.Context.Baggage = make(map[string]string, l)
		for k, v := range parentCtx.Baggage {
			s.raw.Context.Baggage[k] = v
		}
	}
}
