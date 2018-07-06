package trace

import (
	"context"
	"encoding/json"
	"sync"

	"github.com/thundra-io/thundra-lambda-agent-go/plugin"
	"github.com/thundra-io/thundra-lambda-agent-go/otTracer"
	"github.com/opentracing/opentracing-go"
)

type trace struct {
	startTime          int64
	endTime            int64
	duration           int64
	errors             []string
	thrownError        interface{}
	thrownErrorMessage interface{}
	panicInfo          *panicInfo
	errorInfo          *errorInfo
	tracer             opentracing.Tracer
	recorder           *otTracer.InMemorySpanRecorder
}

var invocationCount uint32

// New returns a new trace object.
func New() *trace {
	return &trace{}
}

func NewWithOpentracing() *trace {
	memRecorder := otTracer.NewInMemoryRecorder()
	tracer := otTracer.New(memRecorder)
	opentracing.SetGlobalTracer(tracer)
	return &trace{
		tracer:   tracer,
		recorder: memRecorder,
	}
}

func (trace *trace) GetOpentracingTracer() opentracing.Tracer {
	return trace.tracer
}

func (trace *trace) GetMemRecorder() *otTracer.InMemorySpanRecorder {
	return trace.recorder
}

func (trace *trace) BeforeExecution(ctx context.Context, request json.RawMessage, wg *sync.WaitGroup) {
	cleanBuffer(trace)
	trace.startTime = plugin.GetTimestamp()
	plugin.GenerateNewContextId()
	wg.Done()
}

func (trace *trace) AfterExecution(ctx context.Context, request json.RawMessage, response interface{}, err interface{}) ([]interface{}, string) {
	trace.endTime = plugin.GetTimestamp()
	trace.duration = trace.endTime - trace.startTime

	if err != nil {
		errMessage := plugin.GetErrorMessage(err)
		errType := plugin.GetErrorType(err)

		ei := &errorInfo{
			errMessage,
			errType,
		}

		trace.errorInfo = ei
		trace.thrownError = errType
		trace.thrownErrorMessage = errMessage
		trace.errors = append(trace.errors, errType)
	}

	td := prepareTraceData(request, response, trace)
	var traceArr []interface{}
	traceArr = append(traceArr, td)
	return traceArr, traceDataType
}

func (trace *trace) OnPanic(ctx context.Context, request json.RawMessage, err interface{}, stackTrace []byte) ([]interface{}, string) {
	trace.endTime = plugin.GetTimestamp()
	trace.duration = trace.endTime - trace.startTime

	errMessage := plugin.GetErrorMessage(err)
	errType := plugin.GetErrorType(err)
	pi := &panicInfo{
		errMessage,
		string(stackTrace),
		errType,
	}

	trace.panicInfo = pi
	trace.thrownError = errType
	trace.thrownErrorMessage = plugin.GetErrorMessage(err)
	trace.errors = append(trace.errors, errType)

	td := prepareTraceData(request, nil, trace)
	var traceArr []interface{}
	traceArr = append(traceArr, td)
	return traceArr, traceDataType
}

func cleanBuffer(trace *trace) {
	trace.errors = nil
}
