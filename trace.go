package thundra

import (
	"context"
	"sync"
	"time"
	"github.com/satori/go.uuid"
	"github.com/aws/aws-lambda-go/lambdacontext"
	"os"
	"encoding/json"
	"strings"
	"thundra-agent-go/constants"
	"fmt"
)

type Trace struct {
	startTime          time.Time
	endTime            time.Time
	duration           time.Duration
	errors             []string
	thrownError        interface{}
	thrownErrorMessage interface{}
	panicInfo          *ThundraPanic
	errorInfo          *ThundraError
	collector          *collector
}

type TraceFactory struct{}

type ThundraPanic struct {
	ErrMessage string `json:"errorMessage"`
	StackTrace string `json:"error"`
	ErrType    string `json:"errorType"`
}

type ThundraError struct {
	ErrMessage string `json:"errorMessage"`
	ErrType    string `json:"errorType"`
}

func (t *TraceFactory) Create() Plugin {
	return &Trace{}
}

var invocationCount uint32 = 0
var uniqueId uuid.UUID

type Message struct {
	Data              TraceData `json:"data"`
	Type              string    `json:"type"`
	ApiKey            string    `json:"apiKey"`
	DataFormatVersion string    `json:"dataFormatVersion"`
}

type TraceData struct {
	Id                 string                 `json:"id"`
	ApplicationName    string                 `json:"applicationName"`
	ApplicationId      string                 `json:"applicationId"`
	ApplicationVersion string                 `json:"applicationVersion"`
	ApplicationProfile string                 `json:"applicationProfile"`
	ApplicationType    string                 `json:"applicationType"`
	ContextId          string                 `json:"contextId"`
	ContextName        string                 `json:"contextName"`
	ContextType        string                 `json:"contextType"`
	StartTime          string                 `json:"startTime"`
	EndTime            string                 `json:"endTime"`
	Duration           int64                  `json:"duration"`
	Errors             []string               `json:"errors"`
	ThrownError        interface{}            `json:"thrownError"`
	ThrownErrorMessage interface{}            `json:"thrownErrorMessage"`
	AuditInfo          map[string]interface{} `json:"auditInfo"`
	Properties         map[string]interface{} `json:"properties"`
}

func (trace *Trace) BeforeExecution(ctx context.Context, request interface{}, wg *sync.WaitGroup) {
	defer wg.Done()
	trace.startTime = time.Now()
}

func (trace *Trace) AfterExecution(ctx context.Context, request interface{}, response interface{}, err interface{}, wg *sync.WaitGroup) {
	defer wg.Done()
	trace.endTime = time.Now()
	trace.duration = trace.endTime.Sub(trace.startTime)

	if err != nil {
		trace.errorInfo = &ThundraError{}
		trace.errorInfo.ErrType = getErrorType(err)
		trace.errorInfo.ErrMessage = err.(error).Error()

		trace.thrownError = trace.errorInfo.ErrType
		trace.thrownErrorMessage = trace.errorInfo.ErrMessage
		trace.errors = append(trace.errors, trace.errorInfo.ErrType)
	}

	msg := prepareReport(request, response, err, trace)
	sendReport(trace.collector,msg)
}

func (trace *Trace) OnPanic(ctx context.Context, request json.RawMessage, panic *ThundraPanic, wg *sync.WaitGroup) {
	defer wg.Done()
	trace.endTime = time.Now()
	trace.duration = trace.endTime.Sub(trace.startTime)
	trace.panicInfo = panic

	trace.thrownError = panic.ErrType
	trace.thrownErrorMessage = panic.ErrMessage
	trace.errors = append(trace.errors, panic.ErrType)

	msg := prepareReport(request, nil, nil, trace)
	sendReport(trace.collector,msg)
}

func (trace *Trace)SetCollector(collector *collector) {
	trace.collector = collector
}

func prepareReport(request interface{}, response interface{}, err interface{}, trace *Trace) Message {
	uniqueId = uuid.Must(uuid.NewV4())

	props := prepareProperties(request, response)
	ai := prepareAuditInfo(trace)
	td := prepareTraceData(trace, err, props, ai)
	return prepareMessage(trace, td)
}

func prepareProperties(request interface{}, response interface{}) map[string]interface{} {
	coldStart := "true"
	if invocationCount += 1; invocationCount != 1 {
		coldStart = "false"
	}
	return map[string]interface{}{
		"request":                 request,
		"response":                response,
		"coldStart":               coldStart,
		"functionRegion":          os.Getenv("AWS_DEFAULT_REGION"),
		"functionMemoryLimitInMB": lambdacontext.MemoryLimitInMB,
	}
}

func prepareAuditInfo(trace *Trace) map[string]interface{} {
	var auditErrors []interface{}
	var auditThrownError interface{}

	if trace.panicInfo != nil {
		fmt.Println("PANIC NOT NULL")
		p := *trace.panicInfo
		auditErrors = append(auditErrors, p)
		auditThrownError = p
	} else if trace.errorInfo != nil {
		fmt.Println("ERROR NOT NULL")
		e := *trace.errorInfo
		auditErrors = append(auditErrors, e)
		auditThrownError = e
	}

	return map[string]interface{}{
		"contextName": lambdacontext.FunctionName,
		"id":          uniqueId,
		"openTime":    trace.startTime.Format(constants.TimeFormat),
		"closeTime":   trace.endTime.Format(constants.TimeFormat),
		"errors":      auditErrors,
		"thrownError": auditThrownError,
		//"thrownErrorMessage": trace.thrownErrorMessage,
	}
}

func prepareTraceData(trace *Trace, err interface{}, props map[string]interface{}, auditInfo map[string]interface{}) TraceData {
	appId := splitAppId(lambdacontext.LogStreamName)
	ver := lambdacontext.FunctionVersion

	profile := os.Getenv("thundra_applicationProfile")
	if profile == "" {
		profile = "default"
	}

	return TraceData{
		uniqueId.String(),
		lambdacontext.FunctionName,
		appId,
		ver,
		profile,
		"go",
		uniqueId.String(),
		lambdacontext.FunctionName,
		"ExecutionContext",
		trace.startTime.Format(constants.TimeFormat),
		trace.endTime.Format(constants.TimeFormat),
		convertToMsec(trace.duration), //Convert it to msec
		trace.errors,
		trace.thrownError,
		trace.thrownErrorMessage,
		auditInfo,
		props,
	}
}

func splitAppId(logStreamName string) string {
	s := strings.Split(logStreamName, "]")
	return s[1]
}

func convertToMsec(duration time.Duration) int64 {
	return int64(duration / time.Millisecond)
}

func prepareMessage(trace *Trace, td TraceData) Message {
	return Message{
		td,
		"AuditData",
		ApiKey,
		"1.0",
	}
}
