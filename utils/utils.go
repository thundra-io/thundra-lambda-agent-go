package utils

import (
	"context"
	"fmt"
	"os"
	"reflect"
	"strconv"
	"time"

	opentracing "github.com/opentracing/opentracing-go"
	uuid "github.com/satori/go.uuid"
	"github.com/shirou/gopsutil/process"
	"github.com/thundra-io/thundra-lambda-agent-go/constants"
)

type key struct{}
type eventTypeKey key

// GetTimestamp returns current unix timestamp in msec.
func GetTimestamp() int64 {
	return TimeToMs(time.Now())
}

func TimeToMs(t time.Time) int64 {
	return t.UnixNano() / (int64(time.Millisecond) / int64(time.Nanosecond))
}

func MsToTime(t int64) time.Time {
	return time.Unix(0, t*(int64(time.Millisecond)/int64(time.Nanosecond)))
}

// GenerateNewID generates new uuid.
func GenerateNewID() string {
	return uuid.NewV4().String()
}

// GetThisProcess returns process info about this process.
func GetThisProcess() *process.Process {
	pid := os.Getpid()
	p, err := process.NewProcess(int32(pid))
	if err != nil {
		fmt.Println(err)
		return nil
	}

	return p
}

// GetPid returns pid of this process.
func GetPid() string {
	pid := os.Getpid()
	return strconv.Itoa(pid)
}

// GetErrorType returns type of the error
func GetErrorType(err interface{}) string {
	errorType := reflect.TypeOf(err)
	if errorType.Kind() == reflect.Ptr {
		return errorType.Elem().Name()
	}
	return errorType.Name()
}

// IsTimeout returns whether or not an err is a timeout error
func IsTimeout(err interface{}) bool {
	if err == nil {
		return false
	}
	if GetErrorType(err) == "timeoutError" {
		return true
	}
	return false
}

// GetErrorMessage returns error message
func GetErrorMessage(err interface{}) string {
	e, ok := err.(error)
	if !ok {
		return err.(string)
	}
	return e.Error()
}

// GetEventTypeFromContext returns event type passed in context
func GetEventTypeFromContext(ctx context.Context) interface{} {
	return ctx.Value(eventTypeKey{})
}

// SetEventTypeToContext returns a context with event type value
func SetEventTypeToContext(ctx context.Context, et reflect.Type) context.Context {
	return context.WithValue(ctx, eventTypeKey{}, et)
}

// SetSpanError sets the tags related to the given error to the given span
func SetSpanError(span opentracing.Span, err interface{}) {
	span.SetTag(constants.AwsError, true)
	span.SetTag(constants.AwsErrorKind, GetErrorType(err))
	span.SetTag(constants.AwsErrorMessage, GetErrorMessage(err))
}
