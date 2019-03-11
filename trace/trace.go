package trace

import (
	"context"
	"encoding/json"

	"github.com/thundra-io/thundra-lambda-agent-go/config"

	opentracing "github.com/opentracing/opentracing-go"
	"github.com/thundra-io/thundra-lambda-agent-go/application"
	"github.com/thundra-io/thundra-lambda-agent-go/constants"
	"github.com/thundra-io/thundra-lambda-agent-go/plugin"
	"github.com/thundra-io/thundra-lambda-agent-go/tracer"
	"github.com/thundra-io/thundra-lambda-agent-go/utils"
)

type tracePlugin struct {
	data     *traceData // Not opentracing data just to construct trace plugin data
	rootSpan opentracing.Span
	recorder tracer.SpanRecorder
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
	recorder := tracer.NewInMemoryRecorder()
	tracer := tracer.New(recorder)
	opentracing.SetGlobalTracer(tracer)

	return &tracePlugin{
		recorder: recorder,
	}
}

func (tr *tracePlugin) IsEnabled() bool {
	return !config.TraceDisabled
}

func (tr *tracePlugin) Order() uint8 {
	return pluginOrder
}

// BeforeExecution executes the necessary tasks before the invocation
func (tr *tracePlugin) BeforeExecution(ctx context.Context, request json.RawMessage) context.Context {
	invocationCount++

	startTimeInMs, ctx := plugin.StartTimeFromContext(ctx)
	startTime := utils.MsToTime(startTimeInMs)
	rootSpan, ctx := opentracing.StartSpanFromContext(ctx, application.ApplicationName, opentracing.StartTime(startTime))
	tr.rootSpan = rootSpan

	tr.data = &traceData{
		startTime: startTimeInMs,
	}

	return ctx
}

// AfterExecution executes the necessary tasks after the invocation
func (tr *tracePlugin) AfterExecution(ctx context.Context, request json.RawMessage, response interface{}, err interface{}) ([]plugin.MonitoringDataWrapper, context.Context) {
	finishTime, ctx := plugin.EndTimeFromContext(ctx)
	tr.data.finishTime = finishTime
	tr.data.duration = tr.data.finishTime - tr.data.startTime
	tr.rootSpan.FinishWithOptions(opentracing.FinishOptions{FinishTime: utils.MsToTime(finishTime)})

	// Add root span data
	rawRootSpan, ok := tracer.GetRaw(tr.rootSpan)
	if ok {
		rawRootSpan.ClassName = "AWS-Lambda"
		rawRootSpan.DomainName = "API"
	}

	// Adding tags related to the root span
	tr.rootSpan.SetTag(constants.AwsLambdaName, application.ApplicationName)
	tr.rootSpan.SetTag(constants.AwsLambdaARN, application.GetInvokedFunctionArn(ctx))
	tr.rootSpan.SetTag(constants.AwsRegion, application.FunctionRegion)
	tr.rootSpan.SetTag(constants.AwsLambdaMemoryLimit, application.MemoryLimit)
	tr.rootSpan.SetTag(constants.AwsLambdaLogGroupName, application.LogGroupName)
	tr.rootSpan.SetTag(constants.AwsLambdaLogStreamName, application.LogStreamName)
	tr.rootSpan.SetTag(constants.AwsLambdaInvocationColdStart, invocationCount == 1)
	tr.rootSpan.SetTag(constants.AwsLambdaInvocationTimeout, utils.IsTimeout(err))
	tr.rootSpan.SetTag(constants.AwsLambdaInvocationRequestId, application.GetAwsRequestID(ctx))
	tr.rootSpan.SetTag(constants.AwsLambdaInvocationRequest, request)
	// TODO: Serialize response properly
	tr.rootSpan.SetTag(constants.AwsLambdaInvocationResponse, response)

	if err != nil {
		errMessage := utils.GetErrorMessage(err)
		errType := utils.GetErrorType(err)
		ei := &errorInfo{
			errMessage,
			errType,
		}

		// Add error related tags to the root span
		tr.rootSpan.SetTag(constants.AwsError, true)
		tr.rootSpan.SetTag(constants.AwsErrorKind, errType)
		tr.rootSpan.SetTag(constants.AwsErrorMessage, errMessage)

		tr.data.errorInfo = ei
		tr.data.thrownError = errType
		tr.data.thrownErrorMessage = errMessage
		tr.data.errors = append(tr.data.errors, errType)
	}

	tr.data.timeout = utils.IsTimeout(err)

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

	return traceArr, ctx
}

// Reset clears the recorded data for the next invocation
func (tr *tracePlugin) Reset() {
	tr.recorder.Reset()
}
