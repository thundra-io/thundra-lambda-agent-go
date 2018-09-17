package trace

import (
	"context"
	"encoding/json"
	"sync"

	"github.com/thundra-io/thundra-lambda-agent-go/plugin"
)

type trace struct {
	span *traceSpan
}

// traceSpan collects information related to trace plugin per invocation.
type traceSpan struct {
	startTime          int64
	endTime            int64
	duration           int64
	errors             []string
	thrownError        interface{}
	thrownErrorMessage interface{}
	panicInfo          *panicInfo
	errorInfo          *errorInfo
	timeout            string
}

var invocationCount uint32

// New returns a new trace object.
func New() *trace {
	return new(trace)
}

func (tr *trace) BeforeExecution(ctx context.Context, request json.RawMessage, wg *sync.WaitGroup) {
	tr.span = new(traceSpan)
	tr.span.startTime = plugin.GetTimestamp()
	plugin.GenerateNewContextId()
	wg.Done()
}

func (tr *trace) AfterExecution(ctx context.Context, request json.RawMessage, response interface{}, err interface{}) ([]interface{}, string) {
	tr.span.endTime = plugin.GetTimestamp()
	tr.span.duration = tr.span.endTime - tr.span.startTime

	if err != nil {
		errMessage := plugin.GetErrorMessage(err)
		errType := plugin.GetErrorType(err)

		ei := &errorInfo{
			errMessage,
			errType,
		}

		tr.span.errorInfo = ei
		tr.span.thrownError = errType
		tr.span.thrownErrorMessage = errMessage
		tr.span.errors = append(tr.span.errors, errType)
	}

	tr.span.timeout = isTimeout(err)

	td := tr.prepareTraceData(ctx, request, response)
	tr.span = nil

	var traceArr []interface{}
	traceArr = append(traceArr, td)
	return traceArr, traceDataType
}

func (tr *trace) OnPanic(ctx context.Context, request json.RawMessage, err interface{}, stackTrace []byte) ([]interface{}, string) {
	tr.span.endTime = plugin.GetTimestamp()
	tr.span.duration = tr.span.endTime - tr.span.startTime

	errMessage := plugin.GetErrorMessage(err)
	errType := plugin.GetErrorType(err)
	pi := &panicInfo{
		errMessage,
		string(stackTrace),
		errType,
	}

	tr.span.panicInfo = pi
	tr.span.thrownError = errType
	tr.span.thrownErrorMessage = plugin.GetErrorMessage(err)
	tr.span.errors = append(tr.span.errors, errType)

	// since it is panicked it could not be timed out
	tr.span.timeout = "false"

	td := tr.prepareTraceData(ctx, request, nil)
	tr.span = nil

	var traceArr []interface{}
	traceArr = append(traceArr, td)
	return traceArr, traceDataType
}

// isTimeout returns if the lambda invocation is timed out.
func isTimeout(err interface{}) string {
	if err == nil {
		return "false"
	}
	if plugin.GetErrorType(err) == "timeoutError" {
		return "true"
	}
	return "false"
}
