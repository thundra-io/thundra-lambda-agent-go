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
	"net/http"
	"bytes"
	"io/ioutil"
)

type Trace struct {
	startTime time.Time
	endTime   time.Time
	duration  time.Duration
}

var invocationCount uint32
var id uuid.UUID

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

func prepareReport(request interface{}, response interface{}, error interface{}, trace *Trace) {
	id = uuid.Must(uuid.NewV4())

	props := prepareProperties(request, response)
	ai := prepareAuditInfo(trace, id)
	td := prepareTraceData(trace, error, props, ai)
	msg := prepareMessage(trace, td)

	sendReport(msg)
}

func sendReport(msg Message) {
	b, err := json.Marshal(&msg)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(string(b))

	req, _ := http.NewRequest("POST", url, bytes.NewBuffer(b))
	req.Header.Set("Authorization", "ApiKey "+msg.ApiKey)
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Error:", err)
		panic(err)
	}
	defer resp.Body.Close()

	fmt.Println("response Status:", resp.Status)
	fmt.Println("response Headers:", resp.Header)
	body, _ := ioutil.ReadAll(resp.Body)
	fmt.Println("response Body:", string(body))
}

func prepareAuditInfo(trace *Trace, uuid uuid.UUID) map[string]interface{} {
	return map[string]interface{}{
		"contextName": lambdacontext.FunctionName,
		"id":          uuid,
		"openTime":    trace.startTime.Format(timeFormat),
		"closeTime":   trace.endTime.Format(timeFormat),
	}
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

func prepareTraceData(trace *Trace, err interface{}, props map[string]interface{}, auditInfo map[string]interface{}) TraceData {
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
		id.String(),
		lambdacontext.FunctionName,
		appId,
		ver,
		profile,
		"GO",
		id.String(),
		lambdacontext.FunctionName,
		"ExecutionContext",
		trace.startTime.Format(timeFormat),
		trace.endTime.Format(timeFormat),
		int64(trace.duration / time.Millisecond), //Convert it to msec
		errors,
		auditInfo,
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
