package utils

import (
	"bytes"
	"context"
	"encoding/json"
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

// AttributeValue is a struct for marshalling dynamodb attribute
// value without null fields
type AttributeValue struct {
	B    []byte                      `type:"blob" json:",omitempty"`
	BOOL *bool                       `type:"boolean" json:",omitempty"`
	BS   [][]byte                    `type:"list" json:",omitempty"`
	L    []*AttributeValue           `type:"list" json:",omitempty"`
	M    *map[string]*AttributeValue `type:"map" json:",omitempty"`
	N    *string                     `type:"string" json:",omitempty"`
	NS   []*string                   `type:"list" json:",omitempty"`
	NULL *bool                       `type:"boolean" json:",omitempty"`
	S    *string                     `type:"string" json:"S,omitempty"`
	SS   []*string                   `type:"list" json:",omitempty"`
}

// MarshalJSON implements custom marshaling for AttributeValue
func MarshalJSON(av AttributeValue) ([]byte, error) {
	var buff bytes.Buffer
	var err error
	var b []byte

	if av.B != nil {
		buff.WriteString(`{B: `)
		b = []byte(fmt.Sprintf("%v", string(av.B)))
		buff.Write(b)
	} else if av.BOOL != nil {
		buff.WriteString(`{BOOL: `)
		b, err = json.Marshal(av.BOOL)
		buff.Write(b)
	} else if av.BS != nil {
		buff.WriteString(`{BS: `)
		b, err = json.Marshal(av.BS)
		buff.Write(b)
	} else if av.L != nil {
		buff.WriteString(`{L: `)
		b, err = json.Marshal(av.L)
		buff.Write(b)
	} else if av.M != nil {
		buff.WriteString(`{M: `)
		b, err = json.Marshal(av.M)
		buff.Write(b)
	} else if av.N != nil {
		buff.WriteString(`{N: `)
		b, err = json.Marshal(av.N)
		buff.Write(b)
	} else if av.NS != nil {
		buff.WriteString(`{NS: `)
		b, err = json.Marshal(av.NS)
		buff.Write(b)
	} else if av.NULL != nil {
		buff.WriteString(`{NULL: `)
		b, err = json.Marshal(av.NULL)
		buff.Write(b)
	} else if av.S != nil {
		buff.WriteString(`{S: `)
		b = []byte(fmt.Sprintf("%v", *av.S))
		buff.Write(b)
	} else if av.SS != nil {
		buff.WriteString(`{SS: `)
		b, err = json.Marshal(av.SS)
		buff.Write(b)
	}
	buff.WriteString(`}`)
	return buff.Bytes(), err
}

// AttributeValuetoStr returns string representation of an attribute value
func AttributeValuetoStr(av interface{}) (string, error) {
	attributeValue := AttributeValue{}
	attributeValueJSON, err := json.Marshal(av)
	if err != nil {
		return "", err
	}
	json.Unmarshal(attributeValueJSON, &attributeValue)
	attributeValueBytes, err := MarshalJSON(attributeValue)
	if err != nil {
		return "", err
	}
	return string(attributeValueBytes), nil
}

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

func GetStringFieldFromValue(value reflect.Value, fieldName string) (string, bool) {
	field := value.FieldByName(fieldName)
	if field != (reflect.Value{}) {
		fieldElem := field.Elem()
		if !fieldElem.IsValid() {
			return "", false
		}
		return fieldElem.String(), true
	}
	return "", false
}
