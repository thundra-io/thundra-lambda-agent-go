package trace

import (
	"context"
	"encoding/json"
	"sync"

	"github.com/opentracing/opentracing-go"
	"github.com/thundra-io/thundra-lambda-agent-go/plugin"
	"github.com/thundra-io/thundra-lambda-agent-go/ttracer"
)

type trace struct {
	data     *traceData // Not opentracing data just to construct trace plugin data
	rootSpan opentracing.Span
	recorder ttracer.SpanRecorder
}

// traceData collects information related to trace plugin per invocation.
type traceData struct {
	startTime          int64
	finishTime         int64
	duration           int64
	errors             []string
	thrownError        interface{}
	thrownErrorMessage interface{}
	panicInfo          *panicInfo
	errorInfo          *errorInfo
	timeout            bool
}

var invocationCount uint32

// New returns a new trace object.
func New() *trace {
	recorder := ttracer.NewInMemoryRecorder()
	tracer := ttracer.New(recorder)
	opentracing.SetGlobalTracer(tracer)
	return &trace{
		recorder: recorder,
	}
}

func (tr *trace) BeforeExecution(ctx context.Context, request json.RawMessage, wg *sync.WaitGroup) {
	rootSpan, ctxWithRootSpan := opentracing.StartSpanFromContext(ctx, plugin.FunctionName)
	plugin.CtxWithRootSpan = ctxWithRootSpan
	invocationCount++

	tr.rootSpan = rootSpan
	tr.data = &traceData{
		startTime: plugin.GetTimestamp(),
	}

	wg.Done()
}

func (tr *trace) AfterExecution(ctx context.Context, request json.RawMessage, response interface{}, err interface{}) []plugin.MonitoringDataWrapper {
	tr.rootSpan.Finish()
	tr.data.finishTime = plugin.GetTimestamp()
	tr.data.duration = tr.data.finishTime - tr.data.startTime

	// Adding tags
	tr.rootSpan.SetTag(plugin.AwsLambdaName, plugin.FunctionName)
	tr.rootSpan.SetTag(plugin.AwsLambdaARN, plugin.FunctionARN)
	tr.rootSpan.SetTag(plugin.AwsRegion, plugin.FunctionRegion)

	if err != nil {
		errMessage := plugin.GetErrorMessage(err)
		errType := plugin.GetErrorType(err)

		ei := &errorInfo{
			errMessage,
			errType,
		}

		tr.data.errorInfo = ei
		tr.data.thrownError = errType
		tr.data.thrownErrorMessage = errMessage
		tr.data.errors = append(tr.data.errors, errType)
	}
	tr.data.timeout = isTimeout(err)

	var traceArr []plugin.MonitoringDataWrapper
	td := tr.prepareTraceDataModel(ctx, request, response)
	traceArr = append(traceArr, plugin.WrapMonitoringData(td, traceType))

	spanList := tr.recorder.GetSpans()
	for _, s := range spanList {
		sd := tr.prepareSpanDataModel(ctx, s)
		traceArr = append(traceArr, plugin.WrapMonitoringData(sd, spanType))
	}
	tr.data = nil

	return traceArr
}

func (tr *trace) OnPanic(ctx context.Context, request json.RawMessage, err interface{}, stackTrace []byte) []plugin.MonitoringDataWrapper {
	tr.rootSpan.Finish()
	tr.data.finishTime = plugin.GetTimestamp()
	tr.data.duration = tr.data.finishTime - tr.data.startTime

	errMessage := plugin.GetErrorMessage(err)
	errType := plugin.GetErrorType(err)
	pi := &panicInfo{
		errMessage,
		string(stackTrace),
		errType,
	}

	tr.data.panicInfo = pi
	tr.data.thrownError = errType
	tr.data.thrownErrorMessage = plugin.GetErrorMessage(err)
	tr.data.errors = append(tr.data.errors, errType)
	// Since it is panicked it could not be timed out
	tr.data.timeout = false

	var traceArr []plugin.MonitoringDataWrapper
	td := tr.prepareTraceDataModel(ctx, request, nil)
	traceArr = append(traceArr, plugin.WrapMonitoringData(td, traceType))

	spanList := tr.recorder.GetSpans()
	for _, s := range spanList {
		sd := tr.prepareSpanDataModel(ctx, s)
		traceArr = append(traceArr, plugin.WrapMonitoringData(sd, spanType))
	}

	tr.data = nil

	return traceArr
}

// isTimeout returns if the lambda invocation is timed out.
func isTimeout(err interface{}) bool {
	if err == nil {
		return false
	}
	if plugin.GetErrorType(err) == "timeoutError" {
		return true
	}
	return false
}
