package thundra

import (
	"context"
	"sync"
	"time"
	"github.com/satori/go.uuid"
	"github.com/aws/aws-lambda-go/lambdacontext"
	"os"
	"fmt"
	"encoding/json"
	"strings"
)

type Trace struct {
	startTime          time.Time
	endTime            time.Time
	duration           time.Duration
	errors             []string
	thrownError        interface{}
	thrownErrorMessage interface{}
}

var invocationCount uint32
var uniqueId uuid.UUID

const url = "https://collector.thundra.io/api/audit"

func init() {
	invocationCount = 0
}

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

func (trace *Trace) AfterExecution(ctx context.Context, request interface{}, response interface{}, error interface{}, wg *sync.WaitGroup) {
	defer wg.Done()
	trace.endTime = time.Now()
	trace.duration = trace.endTime.Sub(trace.startTime)
	prepareReport(request, response, error, trace)
}

func (trace *Trace) OnPanic(ctx context.Context, request json.RawMessage, panic *ThundraPanic, wg *sync.WaitGroup) {
	defer wg.Done()
	trace.endTime = time.Now()
	trace.duration = trace.endTime.Sub(trace.startTime)
	preparePanic(trace, panic)

	prepareReport(request, nil, nil, trace)
}

func prepareReport(request interface{}, response interface{}, error interface{}, trace *Trace) {
	uniqueId = uuid.Must(uuid.NewV4())

	props := prepareProperties(request, response)
	ai := prepareAuditInfo(trace)
	td := prepareTraceData(trace, error, props, ai)
	msg := prepareMessage(trace, td)

	sendReport(msg)
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
	return map[string]interface{}{
		"contextName":        lambdacontext.FunctionName,
		"id":                 uniqueId,
		"openTime":           trace.startTime.Format(timeFormat),
		"closeTime":          trace.endTime.Format(timeFormat),
		"errors":             trace.errors,
		"thrownError":        trace.thrownError,
		"thrownErrorMessage": trace.thrownErrorMessage,
	}
}

func preparePanic(trace *Trace, panic *ThundraPanic) {
	var thrownError interface{}
	var thrownErrorMessage interface{}
	if panic != nil {
		trace.errors = append(trace.errors, panic.ErrType) //TODO consider this
		thrownError = panic.ErrType
		thrownErrorMessage = fmt.Sprintf("%v", panic.ErrInfo)
	}
	trace.thrownError = thrownError
	trace.thrownErrorMessage = thrownErrorMessage
}

func prepareTraceData(trace *Trace, err interface{}, props map[string]interface{}, auditInfo map[string]interface{}) TraceData {
	appId := splitAppId(lambdacontext.LogStreamName)
	ver := lambdacontext.FunctionVersion

	if err != nil {
		//TODO consider this, hint of the century consider this
		trace.errors = append(trace.errors, fmt.Sprintf("%+v", err))
	}

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
		"GO",
		uniqueId.String(),
		lambdacontext.FunctionName,
		"ExecutionContext",
		trace.startTime.Format(timeFormat),
		trace.endTime.Format(timeFormat),
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
	key := os.Getenv("thundra_apiKey")

	return Message{
		td,
		"AuditData",
		key,
		"1.0",
	}
}
