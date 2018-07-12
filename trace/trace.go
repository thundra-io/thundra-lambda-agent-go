package trace

import (
	"context"
	"encoding/json"
	"sync"

	"github.com/thundra-io/thundra-lambda-agent-go/plugin"
	"github.com/thundra-io/thundra-lambda-agent-go/thundra_tracer"
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
	recorder           *thundra_tracer.TreeSpanRecorder
}

var invocationCount uint32

// New returns a new trace object with thundra_tracer for opentracing.
// Inorder to use thundra_tracer to manually instrument your code, follow opentracing format.
// If manual instrumentation is not used, collected span data is ignored.
func New() *trace {
	memRecorder := thundra_tracer.NewTreeSpanRecorder()
	tracer := thundra_tracer.New(memRecorder)
	opentracing.SetGlobalTracer(tracer)
	return &trace{
		recorder: memRecorder,
	}
}

// GetRecorder returns the TreeSpanRecorder
func (trace *trace) GetRecorder() *thundra_tracer.TreeSpanRecorder {
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
