package trace

import (
	"context"
	"encoding/json"
	"sync"

	"github.com/thundra-io/thundra-lambda-agent-go/plugin"
)

type trace struct {
	startTime          int64
	endTime            int64
	duration           int64
	errors             []string
	thrownError        interface{}
	thrownErrorMessage interface{}
	panicInfo          *panicInfo
	errorInfo          *errorInfo
}

var invocationCount uint32

// New returns a new trace object.
func New() *trace {
	return &trace{}
}

func (tr *trace) BeforeExecution(ctx context.Context, request json.RawMessage, wg *sync.WaitGroup) {
	cleanBuffer(tr)
	tr.startTime = plugin.GetTimestamp()
	plugin.GenerateNewContextId()
	wg.Done()
}

func (tr *trace) AfterExecution(ctx context.Context, request json.RawMessage, response interface{}, err interface{}) ([]interface{}, string) {
	tr.endTime = plugin.GetTimestamp()
	tr.duration = tr.endTime - tr.startTime

	if err != nil {
		errMessage := plugin.GetErrorMessage(err)
		errType := plugin.GetErrorType(err)

		ei := &errorInfo{
			errMessage,
			errType,
		}

		tr.errorInfo = ei
		tr.thrownError = errType
		tr.thrownErrorMessage = errMessage
		tr.errors = append(tr.errors, errType)
	}

	td := tr.prepareTraceData(ctx, request, response)
	var traceArr []interface{}
	traceArr = append(traceArr, td)
	return traceArr, traceDataType
}

func (tr *trace) OnPanic(ctx context.Context, request json.RawMessage, err interface{}, stackTrace []byte) ([]interface{}, string) {
	tr.endTime = plugin.GetTimestamp()
	tr.duration = tr.endTime - tr.startTime

	errMessage := plugin.GetErrorMessage(err)
	errType := plugin.GetErrorType(err)
	pi := &panicInfo{
		errMessage,
		string(stackTrace),
		errType,
	}

	tr.panicInfo = pi
	tr.thrownError = errType
	tr.thrownErrorMessage = plugin.GetErrorMessage(err)
	tr.errors = append(tr.errors, errType)

	td := tr.prepareTraceData(ctx, request, nil)
	var traceArr []interface{}
	traceArr = append(traceArr, td)
	return traceArr, traceDataType
}

func cleanBuffer(trace *trace) {
	trace.errors = nil
}
