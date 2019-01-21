package trace

import (
	"context"
	"encoding/json"
	"os"

	opentracing "github.com/opentracing/opentracing-go"
	"github.com/thundra-io/thundra-lambda-agent-go/plugin"
	"github.com/thundra-io/thundra-lambda-agent-go/ttracer"
)

type tracePlugin struct {
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
func New() *tracePlugin {
	recorder := ttracer.NewInMemoryRecorder()
	tracer := ttracer.New(recorder)
	opentracing.SetGlobalTracer(tracer)
	return &tracePlugin{
		recorder: recorder,
	}
}

func (tr *tracePlugin) IsEnabled() bool {
	if os.Getenv(plugin.ThundraDisableTrace) == "true" {
		return false
	}

	return true
}

// BeforeExecution executes the necessary tasks before the invocation
func (tr *tracePlugin) BeforeExecution(ctx context.Context, request json.RawMessage) context.Context {
	rootSpan, ctxWithRootSpan := opentracing.StartSpanFromContext(ctx, plugin.FunctionName)
	invocationCount++

	tr.rootSpan = rootSpan
	tr.data = &traceData{
		startTime: plugin.GetTimestamp(),
	}

	return ctxWithRootSpan
}

// AfterExecution executes the necessary tasks after the invocation
func (tr *tracePlugin) AfterExecution(ctx context.Context, request json.RawMessage, response interface{}, err interface{}) []plugin.MonitoringDataWrapper {
	tr.rootSpan.Finish()
	tr.data.finishTime = plugin.GetTimestamp()
	tr.data.duration = tr.data.finishTime - tr.data.startTime

	// Add root span data
	rawRootSpan, ok := ttracer.GetRaw(tr.rootSpan)
	if ok {
		rawRootSpan.ClassName = "AWS-Lambda"
		rawRootSpan.DomainName = "API"
	}

	// Adding tags related to the root span
	tr.rootSpan.SetTag(plugin.AwsLambdaName, plugin.FunctionName)
	tr.rootSpan.SetTag(plugin.AwsLambdaARN, plugin.GetInvokedFunctionArn(ctx))
	tr.rootSpan.SetTag(plugin.AwsRegion, plugin.FunctionRegion)
	tr.rootSpan.SetTag(plugin.AwsLambdaMemoryLimit, plugin.MemoryLimit)
	tr.rootSpan.SetTag(plugin.AwsLambdaLogGroupName, plugin.LogGroupName)
	tr.rootSpan.SetTag(plugin.AwsLambdaLogStreamName, plugin.LogStreamName)
	tr.rootSpan.SetTag(plugin.AwsLambdaInvocationColdStart, invocationCount == 1)
	tr.rootSpan.SetTag(plugin.AwsLambdaInvocationTimeout, plugin.IsTimeout(err))
	tr.rootSpan.SetTag(plugin.AwsLambdaInvocationRequestId, plugin.GetAwsRequestID(ctx))
	tr.rootSpan.SetTag(plugin.AwsLambdaInvocationRequest, request)
	tr.rootSpan.SetTag(plugin.AwsLambdaInvocationResponse, response)

	if err != nil {
		errMessage := plugin.GetErrorMessage(err)
		errType := plugin.GetErrorType(err)
		ei := &errorInfo{
			errMessage,
			errType,
		}

		// Add error related tags to the root span
		tr.rootSpan.SetTag(plugin.AwsError, true)
		tr.rootSpan.SetTag(plugin.AwsErrorKind, errType)
		tr.rootSpan.SetTag(plugin.AwsErrorMessage, errMessage)

		tr.data.errorInfo = ei
		tr.data.thrownError = errType
		tr.data.thrownErrorMessage = errMessage
		tr.data.errors = append(tr.data.errors, errType)
	}

	tr.data.timeout = plugin.IsTimeout(err)

	// Prepare report data
	var traceArr []plugin.MonitoringDataWrapper
	td := tr.prepareTraceDataModel(ctx, request, response)
	traceArr = append(traceArr, plugin.WrapMonitoringData(td, traceType))

	spanList := tr.recorder.GetSpans()
	for _, s := range spanList {
		sd := tr.prepareSpanDataModel(ctx, s)
		traceArr = append(traceArr, plugin.WrapMonitoringData(sd, spanType))
	}

	// Clear trace plugin data for next invocation
	tr.Reset()

	return traceArr
}

// Reset clears the recorded data for the next invocation
func (tr *tracePlugin) Reset() {
	tr.data = nil
	tr.recorder.Reset()
}
