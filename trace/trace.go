package trace

import (
	"context"
	"sync"
	"time"
	"os"
	"encoding/json"
	"strings"
	"fmt"
	"reflect"

	"github.com/satori/go.uuid"
	"github.com/aws/aws-lambda-go/lambdacontext"
)

type Trace struct {
	startTime          time.Time
	endTime            time.Time
	duration           time.Duration
	errors             []string
	thrownError        interface{}
	thrownErrorMessage interface{}
	panicInfo          *panicInfo
	errorInfo          *errorInfo
}

var invocationCount uint32 = 0
var uniqueId uuid.UUID

type traceData struct {
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

func (trace *Trace) BeforeExecution(ctx context.Context, request json.RawMessage, wg *sync.WaitGroup) {
	trace.startTime = time.Now().Round(time.Millisecond)
	cleanBuffer(trace)
	wg.Done()
}

func (trace *Trace) AfterExecution(ctx context.Context, request json.RawMessage, response interface{}, err interface{}, wg *sync.WaitGroup) (interface{}, string) {
	defer wg.Done()
	trace.endTime = time.Now().Round(time.Millisecond)
	trace.duration = trace.endTime.Sub(trace.startTime)

	if err != nil {
		errMessage := getErrorMessage(err)
		errType := getErrorType(err)

		ei := &errorInfo{
			errMessage,
			errType,
		}

		trace.errorInfo = ei
		trace.thrownError = errType
		trace.thrownErrorMessage = errMessage
		trace.errors = append(trace.errors, errType)
	}

	td := prepareReport(request, response, err, trace)
	return td, TraceDataType
}

func (trace *Trace) OnPanic(ctx context.Context, request json.RawMessage, err interface{}, stackTrace []byte, wg *sync.WaitGroup) (interface{}, string) {
	defer wg.Done()
	trace.endTime = time.Now()
	trace.duration = trace.endTime.Sub(trace.startTime)

	errMessage := err.(error).Error()
	errType := getErrorType(err)
	pi := &panicInfo{
		errMessage,
		string(stackTrace),
		errType,
	}

	trace.panicInfo = pi
	trace.thrownError = errType
	trace.thrownErrorMessage = getErrorMessage(err)
	trace.errors = append(trace.errors, errType)

	td := prepareReport(request, nil, nil, trace)
	return td, TraceDataType
}

func cleanBuffer(trace *Trace) {
	trace.errors = nil
}

func getErrorType(err interface{}) string {
	errorType := reflect.TypeOf(err)
	if errorType.Kind() == reflect.Ptr {
		return errorType.Elem().Name()
	}
	return errorType.Name()
}

func getErrorMessage(err interface{}) string {
	return err.(error).Error()
}

func prepareReport(request json.RawMessage, response interface{}, err interface{}, trace *Trace) interface{} {
	uniqueId = uuid.Must(uuid.NewV4())

	props := prepareProperties(request, response)
	ai := prepareAuditInfo(trace)
	td := prepareTraceData(trace, err, props, ai)
	return td
}

func prepareProperties(request json.RawMessage, response interface{}) map[string]interface{} {
	coldStart := "true"
	if invocationCount += 1; invocationCount != 1 {
		coldStart = "false"
	}
	return map[string]interface{}{
		auditInfoPropertiesRequest:             string(request),
		auditInfoPropertiesResponse:            response,
		auditInfoPropertiesColdStart:           coldStart,
		auditInfoPropertiesFunctionRegion:      os.Getenv(awsDefaultRegion),
		auditInfoPropertiesFunctionMemoryLimit: lambdacontext.MemoryLimitInMB,
	}
}

func prepareAuditInfo(trace *Trace) map[string]interface{} {
	var auditErrors []interface{}
	var auditThrownError interface{}

	if trace.panicInfo != nil {
		fmt.Println("Panic is not null")
		p := *trace.panicInfo
		auditErrors = append(auditErrors, p)
		auditThrownError = p
	} else if trace.errorInfo != nil {
		fmt.Println("Error is not null")
		e := *trace.errorInfo
		auditErrors = append(auditErrors, e)
		auditThrownError = e
	}

	return map[string]interface{}{
		auditInfoContextName: lambdacontext.FunctionName,
		auditInfoId:          uniqueId,
		auditInfoOpenTime:    trace.startTime.Format(timeFormat),
		auditInfoCloseTime:   trace.endTime.Format(timeFormat),
		auditInfoErrors:      auditErrors,
		auditInfoThrownError: auditThrownError,
		//"thrownErrorMessage": trace.thrownErrorMessage,
	}
}

func prepareTraceData(trace *Trace, err interface{}, props map[string]interface{}, auditInfo map[string]interface{}) traceData {
	appId := splitAppId(lambdacontext.LogStreamName)
	ver := lambdacontext.FunctionVersion

	profile := os.Getenv(thundraApplicationProfile)
	if profile == "" {
		profile = defaultProfile
	}

	return traceData{
		uniqueId.String(),
		lambdacontext.FunctionName,
		appId,
		ver,
		profile,
		applicationType,
		uniqueId.String(),
		lambdacontext.FunctionName,
		executionContext,
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
	if len(s) > 1 {
		return s[1]
	} else {
		return ""
	}
}

func convertToMsec(duration time.Duration) int64 {
	return int64(duration / time.Millisecond)
}
