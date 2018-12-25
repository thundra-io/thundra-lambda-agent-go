package ttracer

import (
	ot "github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/log"
	"github.com/thundra-io/thundra-lambda-agent-go/plugin"
	"sync"
	"time"
)

type spanImpl struct {
	tracer     *tracerImpl
	sync.Mutex // protects the fields below
	raw        RawSpan
	// The number of logs dropped because of MaxLogsPerSpan.
	numDroppedLogs int
}

func newSpan(operationName string, tracer *tracerImpl, sso []ot.StartSpanOption) *spanImpl {
	opts := newStartSpanOptions(sso)

	// Start time.
	startTimestamp := plugin.GetTimestamp()

	// Build the new span. This is the only allocation: We'll return this as
	// an opentracing.Span.
	sp := &spanImpl{}

	// Look for a parent in the list of References.
ReferencesLoop:
	for _, ref := range opts.Options.References {
		switch ref.Type {
		case ot.ChildOfRef,
			ot.FollowsFromRef:
			//TODO traverse map, find parent from parentId, then add as a child
			refCtx := ref.ReferencedContext.(SpanContext)
			sp.raw.Context.TraceID = refCtx.TraceID
			sp.raw.ParentSpanID = refCtx.SpanID
			sp.raw.Context.SpanID = plugin.GenerateNewId()

			if l := len(refCtx.Baggage); l > 0 {
				sp.raw.Context.Baggage = make(map[string]string, l)
				for k, v := range refCtx.Baggage {
					sp.raw.Context.Baggage[k] = v
				}
			}
			break ReferencesLoop
		}
	}

	if sp.raw.Context.TraceID == "" {
		// TraceID not set by parent reference or explicitly
		sp.raw.Context.TraceID, sp.raw.Context.SpanID = plugin.Generate2NewId()
	} else if sp.raw.Context.SpanID == "" {
		// TraceID set but SpanID not set
		sp.raw.Context.SpanID = plugin.GenerateNewId()
	}

	sp.tracer = tracer
	sp.raw.Operation = operationName
	sp.raw.StartTimestamp = startTimestamp
	sp.raw.Duration = -1
	sp.raw.operationType = opts.operationType
	sp.raw.operationGroup = opts.operationGroup
	sp.raw.Tags = opts.Options.Tags

	sp.tracer.opts.Recorder.RecordSpanStarted(&sp.raw)

	return sp
}

// Finish is the last call that is made
func (s *spanImpl) Finish() {
	s.FinishWithOptions(ot.FinishOptions{})
}

// FinishWithOptions finishes span and adds the given options to it
func (s *spanImpl) FinishWithOptions(opts ot.FinishOptions) {
	finishTimestamp := plugin.GetTimestamp()
	duration := finishTimestamp - s.StartTimestamp()

	s.Lock()
	defer s.Unlock()

	for _, lr := range opts.LogRecords {
		s.appendLog(lr)
	}
	for _, ld := range opts.BulkLogData {
		s.appendLog(ld.ToLogRecord())
	}

	if s.numDroppedLogs > 0 {
		// We dropped some log events, which means that we used part of Logs as a
		// circular buffer (see appendLog). De-circularize it.
		numOld := (len(s.raw.Logs) - 1) / 2
		numNew := len(s.raw.Logs) - numOld
		rotateLogBuffer(s.raw.Logs[numOld:], s.numDroppedLogs%numNew)

		// Replace the log in the middle (the oldest "new" log) with information
		// about the dropped logs. This means that we are effectively dropping one
		// more "new" log.
		numDropped := s.numDroppedLogs + 1
		s.raw.Logs[numOld] = ot.LogRecord{
			// Keep the timestamp of the last dropped event.
			Timestamp: s.raw.Logs[numOld].Timestamp,
			Fields: []log.Field{
				log.String("event", "dropped Span logs"),
				log.Int("dropped_log_count", numDropped),
				log.String("component", "basictracer"),
			},
		}
	}

	s.raw.Duration = duration
	s.raw.EndTimestamp = finishTimestamp

	s.tracer.opts.Recorder.RecordSpanEnded()
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
	if s.tracer.opts.DropAllLogs {
		return
	}

	if ld.Timestamp.IsZero() {
		ld.Timestamp = time.Now()
	}

	s.appendLog(ld.ToLogRecord())
}

// rotateLogBuffer rotates the records in the buffer: records 0 to pos-1 move at
// the end (i.e. pos circular left shifts).
func rotateLogBuffer(buf []ot.LogRecord, pos int) {
	// This algorithm is described in:
	//    http://www.cplusplus.com/reference/algorithm/rotate
	for first, middle, next := 0, pos, pos; first != middle; {
		buf[first], buf[next] = buf[next], buf[first]
		first++
		next++
		if next == len(buf) {
			next = middle
		} else if first == middle {
			middle = next
		}
	}
}

// SetOperationName sets operation name.
func (s *spanImpl) SetOperationName(operationName string) ot.Span {
	s.Lock()
	defer s.Unlock()
	s.raw.Operation = operationName
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
	maxLogs := s.tracer.opts.MaxLogsPerSpan
	if maxLogs == 0 || len(s.raw.Logs) < maxLogs {
		s.raw.Logs = append(s.raw.Logs, lr)
		return
	}

	// We have too many logs. We don't touch the first numOld logs; we treat the
	// rest as a circular buffer and overwrite the oldest log among those.
	numOld := (maxLogs - 1) / 2
	numNew := maxLogs - numOld
	s.raw.Logs[numOld+s.numDroppedLogs%numNew] = lr
	s.numDroppedLogs++
}

// LogFields parses parameter fields sequentially, as first one is the key and the second is it's value.
func (s *spanImpl) LogFields(fields ...log.Field) {
	lr := ot.LogRecord{
		Fields: fields,
	}
	s.Lock()
	defer s.Unlock()
	if s.tracer.opts.DropAllLogs {
		return
	}
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
func (s *spanImpl) Operation() string {
	return s.raw.Operation
}

// StartTimestamp returns StartTimestamp
func (s *spanImpl) StartTimestamp() int64 {
	return s.raw.StartTimestamp
}
