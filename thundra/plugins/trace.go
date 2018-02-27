package plugins

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
	startTime time.Time
	endTime   time.Time
	duration  time.Duration
}

var invocationCount uint32

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
	Duration           int64                  `json:"duration"`
	StartTime          string                 `json:"startTime"`
	EndTime            string                 `json:"endTime"`
	Errors             []string               `json:"errors"`
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
	sendReport(request, response, error, trace)
}

func sendReport(request interface{}, response interface{}, error interface{}, trace *Trace) {
	props := prepareProperties(request, response)
	td := prepareTraceData(trace, error, props)
	msg := prepareMessage(trace, td)

	b, err := json.Marshal(&msg)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(string(b))
}

func prepareProperties(request interface{}, response interface{}) map[string]interface{} {
	coldStart := true
	if invocationCount += 1; invocationCount != 1 {
		coldStart = false
	}
	props := map[string]interface{}{
		"request":                 request,
		"response":                response,
		"coldStart":               coldStart,
		"functionRegion":          os.Getenv("AWS_DEFAULT_REGION"),
		"functionMemoryLimitInMB": lambdacontext.MemoryLimitInMB,
	}

	return props;
}

func prepareTraceData(trace *Trace, err interface{}, props map[string]interface{}) TraceData {
	u1 := uuid.Must(uuid.NewV4())
	appName := lambdacontext.FunctionName
	appId := splitAppId(lambdacontext.LogStreamName)
	ver := lambdacontext.FunctionVersion

	profile := os.Getenv("thundra_applicationProfile")
	if profile == "" {
		profile = "default"
	}

	errors := []string{}
	if err != nil {
		errors = append(errors, err.(string))
	}

	return TraceData{
		u1.String(),
		appName,
		appId,
		ver,
		profile,
		"GO",
		int64(trace.duration / time.Millisecond), //Convert it to msec
		trace.startTime.Format(timeFormat),
		trace.endTime.Format(timeFormat),
		errors,
		props,
	}
}

func splitAppId(logStreamName string) string {
	s := strings.Split(logStreamName, "]")
	return s[1]
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
