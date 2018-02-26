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

var invocationCount uint32

type Trace struct {
	startTime time.Time
	endTime   time.Time
	duration  time.Duration
}

func init() {
	invocationCount = 0
}

type Message struct {
	Data              TraceData              `json:"data"`
	Type              string                 `json:"type"`
	ApiKey            string                 `json:"apiKey"`
	DataFormatVersion string                 `json:"dataFormatVersion"`
	Duration          time.Duration          `json:"duration"`
	StartTime         time.Time              `json:"startTime"` //TODO “yyyy-MM-dd HH:mm:ss.SSS Z” format
	EndTime           time.Time              `json:"endTime"`   //TODO “yyyy-MM-dd HH:mm:ss.SSS Z” format
	Errors            []string               `json:"errors"`
	Properties        map[string]interface{} `json:"properties"`
}

type TraceData struct {
	Id                 string `json:"id"`
	ApplicationName    string `json:"applicationName"`
	ApplicationId      string `json:"applicationId"`
	ApplicationVersion string `json:"applicationVersion"`
	ApplicationProfile string `json:"applicationProfile"`
	ApplicationType    string `json:"applicationType"`
}

func (trace *Trace) BeforeExecution(ctx context.Context, request interface{}, wg *sync.WaitGroup) {
	defer wg.Done()
	trace.startTime = time.Now()
}

func (trace *Trace) AfterExecution(ctx context.Context, request interface{}, response interface{}, error interface{}, wg *sync.WaitGroup) {
	defer wg.Done()
	trace.duration = time.Since(trace.startTime)
	sendReport(request, response, error, trace)
}

func sendReport(request interface{}, response interface{}, error interface{}, trace *Trace) {
	td := prepareTraceData()
	props := prepareProperties(request, response)
	msg := prepareMessage(trace, td, error, props)
	fmt.Println("Message:\n", msg)
}

func prepareProperties(request interface{}, response interface{}) map[string]interface{} {
	coldStart := true
	if invocationCount != 0 {
		coldStart = false
	}
	req, _ := json.Marshal(&request)
	props := map[string]interface{}{
		"request":                 string(req),
		"response":                response,
		"coldStart":               coldStart,
		"functionRegion":          os.Getenv("AWS_DEFAULT_REGION"),
		"FunctionMemoryLimitInMB": lambdacontext.MemoryLimitInMB,
	}

	return props;
}

func prepareMessage(trace *Trace, td TraceData, error interface{}, props map[string]interface{}) Message {
	key := os.Getenv("thundra_apiKey")
	errors := []string{}
	if error != nil {
		errors = append(errors, error.(string))
	}
	return Message{
		td,
		"AuditData",
		key,
		"1.0",
		trace.duration,
		trace.startTime,
		trace.endTime,
		errors,
		props}
}

func splitAppId(logStreamName string) string {
	s := strings.Split(logStreamName, "]")
	return s[1]
}

func prepareTraceData() TraceData {
	u1 := uuid.Must(uuid.NewV4())
	appName := lambdacontext.FunctionName
	appId := splitAppId(lambdacontext.LogStreamName)
	ver := lambdacontext.FunctionVersion
	profile := os.Getenv("thundra_applicationProfile")

	return TraceData{u1.String(), appName, appId, ver, profile, "GO"}
}
