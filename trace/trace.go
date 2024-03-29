package trace

import (
	"context"
	"encoding/json"
	logger "log"
	"sync"

	"github.com/aws/aws-lambda-go/events"
	"github.com/thundra-io/thundra-lambda-agent-go/v2/config"

	opentracing "github.com/opentracing/opentracing-go"
	"github.com/thundra-io/thundra-lambda-agent-go/v2/application"
	"github.com/thundra-io/thundra-lambda-agent-go/v2/constants"
	"github.com/thundra-io/thundra-lambda-agent-go/v2/plugin"
	"github.com/thundra-io/thundra-lambda-agent-go/v2/tracer"
	"github.com/thundra-io/thundra-lambda-agent-go/v2/utils"
)

type tracePlugin struct {
	Data     *Data // Not opentracing data just to construct trace plugin data
	RootSpan opentracing.Span
	Recorder tracer.SpanRecorder
}

// Data collects information related to trace plugin per invocation.
type Data struct {
	StartTime          int64
	FinishTime         int64
	Duration           int64
	Errors             []string
	ThrownError        interface{}
	ThrownErrorMessage interface{}
	PanicInfo          *panicInfo
	ErrorInfo          *errorInfo
	Timeout            bool
}

var invocationCount uint32

var lock = &sync.Mutex{}
var instance *tracePlugin

// New returns a new trace object.
func New() *tracePlugin {
	recorder := tracer.NewInMemoryRecorder()
	tracer := tracer.New(recorder)
	opentracing.SetGlobalTracer(tracer)

	return &tracePlugin{
		Recorder: recorder,
	}
}

// GetInstance returns the tracePlugin instance existing or creates
// a new instance and then returns it
func GetInstance() *tracePlugin {
	lock.Lock()
	defer lock.Unlock()

	if instance == nil {
		instance = New()
	}
	return instance
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
	tr.RootSpan = rootSpan

	tr.Data = &Data{
		StartTime: startTimeInMs,
	}

	// Add root span data
	rawRootSpan, ok := tracer.GetRaw(tr.RootSpan)
	if ok {
		rawRootSpan.ClassName = constants.AwsLambdaApplicationClass
		rawRootSpan.DomainName = constants.AwsLambdaApplicationDomain
	}

	// Adding tags related to the root span
	tr.RootSpan.SetTag(constants.AwsLambdaName, application.ApplicationName)
	tr.RootSpan.SetTag(constants.AwsLambdaARN, application.GetInvokedFunctionArn(ctx))
	tr.RootSpan.SetTag(constants.AwsRegion, application.FunctionRegion)
	tr.RootSpan.SetTag(constants.AwsLambdaMemoryLimit, application.MemoryLimit)
	tr.RootSpan.SetTag(constants.AwsLambdaLogGroupName, application.LogGroupName)
	tr.RootSpan.SetTag(constants.AwsLambdaLogStreamName, application.LogStreamName)
	tr.RootSpan.SetTag(constants.AwsLambdaInvocationColdStart, invocationCount == 1)
	tr.RootSpan.SetTag(constants.AwsLambdaInvocationRequestId, application.GetAwsRequestID(ctx))

	tracer.OnSpanStarted(tr.RootSpan)

	return ctx
}

// AfterExecution executes the necessary tasks after the invocation
func (tr *tracePlugin) AfterExecution(ctx context.Context, request json.RawMessage, response interface{}, err interface{}) ([]plugin.MonitoringDataWrapper, context.Context) {
	finishTime, ctx := plugin.EndTimeFromContext(ctx)
	tr.Data.FinishTime = finishTime
	tr.Data.Duration = tr.Data.FinishTime - tr.Data.StartTime

	tr.finishRootSpan()

	tr.RootSpan.SetTag(constants.AwsLambdaInvocationTimeout, utils.IsTimeout(err))

	// Disable request data sending for cloudwatchlog, firehose and kinesis if not
	// enabled by configuration because requests can get too big for these
	enableRequestData := true
	if (config.TraceRequestDisabled) ||
		(plugin.TriggerClassName == constants.ClassNames["KINESIS"] && !config.TraceKinesisRequestEnabled) ||
		(plugin.TriggerClassName == constants.ClassNames["FIREHOSE"] && !config.TraceFirehoseRequestEnabled) ||
		(plugin.TriggerClassName == constants.ClassNames["CLOUDWATCHLOG"] && !config.TraceCloudwatchlogRequestEnabled) {
		enableRequestData = false
	}
	if enableRequestData {
		tr.RootSpan.SetTag(constants.AwsLambdaInvocationRequest, request)
	}

	if plugin.TriggerClassName == constants.ClassNames["APIGATEWAY"] {
		if response != nil {
			responseInterface, ok := (response).(*interface{})
			if ok {
				APIGwResponse, ok := (*responseInterface).(events.APIGatewayProxyResponse)
				if ok {
					if APIGwResponse.Headers == nil {
						APIGwResponse.Headers = make(map[string]string)
					}
					var APIGWRequest events.APIGatewayProxyRequest
					json.Unmarshal(request, &APIGWRequest)
					APIGwResponse.Headers[constants.AwsLambdaTriggerResourceName] = APIGWRequest.Resource

					*responseInterface = APIGwResponse
				}
			}

		}
	}

	if !config.TraceResponseDisabled {
		// TODO: Serialize response properly
		tr.RootSpan.SetTag(constants.AwsLambdaInvocationResponse, response)
	}

	if err != nil {
		errMessage := utils.GetErrorMessage(err)
		errType := utils.GetErrorType(err)
		ei := &errorInfo{
			errMessage,
			errType,
		}

		// Add error related tags to the root span
		tr.RootSpan.SetTag(constants.AwsError, true)
		tr.RootSpan.SetTag(constants.AwsErrorKind, errType)
		tr.RootSpan.SetTag(constants.AwsErrorMessage, errMessage)

		utils.SetSpanError(tr.RootSpan, err)

		tr.Data.ErrorInfo = ei
		tr.Data.ThrownError = errType
		tr.Data.ThrownErrorMessage = errMessage
		tr.Data.Errors = append(tr.Data.Errors, errType)
	}

	tr.Data.Timeout = utils.IsTimeout(err)

	spanList := tr.Recorder.GetSpans()

	sampled := true
	sampler := GetSampler()
	if len(spanList) > 0 && sampler != nil {
		sampled = sampler.IsSampled(spanList[0])
	}
	// Prepare report data
	var traceArr []plugin.MonitoringDataWrapper
	if sampled {
		for _, s := range spanList {
			sd := tr.prepareSpanDataModel(ctx, s)
			traceArr = append(traceArr, plugin.WrapMonitoringData(sd, spanType))
		}
	}

	// Clear trace plugin data for next invocation
	tr.Reset()

	return traceArr, ctx
}

func (tr *tracePlugin) finishRootSpan() {
	defer func() {
		if r := recover(); r != nil {
			logger.Println("Error while finishing the root span:", r)
		}
	}()
	tr.RootSpan.FinishWithOptions(opentracing.FinishOptions{FinishTime: utils.MsToTime(tr.Data.FinishTime)})
}

// Reset clears the recorded data for the next invocation
func (tr *tracePlugin) Reset() {
	tr.Recorder.Reset()
}
