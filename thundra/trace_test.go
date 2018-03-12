package thundra

import (
	"testing"
	"context"
	"fmt"
	"encoding/json"
	"os"
	"time"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/aws/aws-lambda-go/lambdacontext"
	"thundra-agent-go/trace"
	"thundra-agent-go/plugin"
	"errors"
)

const (
	duration           = 500
	functionName       = "TestFunctionName"
	memoryLimit        = 512
	functionVersion    = "$Version"
	applicationProfile = "TestProfile"
	TestApiKey         = "TestApiKey"
	defaultRegion      = "TestRegion"
	generatedError     = "Generated Error"
	errorType          = "errorString"
	generatedPanic     = "Generated Panic"
)

var coldStart = "true"

type expectedPanic struct {
	val   string
	err   error
	panic error
}

func TestTrace(t *testing.T) {
	hello := func(s string) string {
		time.Sleep(time.Millisecond * duration)
		return fmt.Sprintf("%s works!", s)
	}

	testCases := []struct {
		name     string
		input    string
		expected expected
		handler  interface{}
	}{
		{
			input:    `"Thundra"`,
			expected: expected{"Thundra works!", nil},
			handler: func(name string) (string, error) {
				return hello(name), nil
			},
		},
		{
			input:    `"Thundra"`,
			expected: expected{"Thundra works!", nil},
			handler: func(name string) (string, error) {
				return hello(name), nil
			},
		},
		{
			input:    `"Thundra"`,
			expected: expected{"Thundra works!", nil},
			handler: func(ctx context.Context, name string) (string, error) {
				return hello(name), nil
			},
		},
		{
			input:    `"Thundra"`,
			expected: expected{"", errors.New(generatedError)},
			handler: func() error {
				time.Sleep(time.Millisecond * duration)
				return errors.New(generatedError)
			},
		},
		{
			input:    `"Thundra"`,
			expected: expected{"", errors.New(generatedError)},
			handler: func() (interface{}, error) {
				time.Sleep(time.Millisecond * duration)
				return nil, errors.New(generatedError)
			},
		},
		{
			input:    `"Thundra"`,
			expected: expected{"", errors.New(generatedError)},
			handler: func(e interface{}) (interface{}, error) {
				time.Sleep(time.Millisecond * duration)
				return nil, errors.New(generatedError)
			},
		},
		{
			input:    `"Thundra"`,
			expected: expected{"", errors.New(generatedError)},
			handler: func(ctx context.Context, e interface{}) (interface{}, error) {
				time.Sleep(time.Millisecond * duration)
				return nil, errors.New(generatedError)
			},
		},
	}
	for i, testCase := range testCases {
		t.Run(fmt.Sprintf("testCase[%d] %s", i, testCase.name), func(t *testing.T) {
			prepareEnvironment()

			r := new(MockReporter)
			r.On("report").Return()
			r.On("clear").Return()
			r.On("collect", mock.Anything).Return()

			tf := trace.TraceFactory{}
			tr := tf.Create()
			th := NewBuilder().AddPlugin(tr).SetReporter(r).Build()
			lambdaHandler := WrapLambdaHandler(testCase.handler, th)

			response, err, invocationStartTime, invocationEndTime := lambdaHandler.InvokeWithoutSerialization(context.TODO(), []byte(testCase.input))

			//Monitor Data
			msg, ok := r.messageQueue[0].(plugin.Message)
			if !ok {
				fmt.Println("Collector message can't be casted to pluginMessage")
			}
			assert.Equal(t, trace.TraceDataType, msg.Type)
			assert.Equal(t, TestApiKey, msg.ApiKey)
			assert.Equal(t, dataFormatVersion, msg.DataFormatVersion)

			//Trace Data
			td, ok := msg.Data.(trace.TraceData)
			if !ok {
				fmt.Println("Can not convert to trace data")
			}
			assert.NotNil(t, td.Id)
			assert.Equal(t, functionName, td.ApplicationName)
			assert.Equal(t, "", td.ApplicationId)
			assert.Equal(t, functionVersion, td.ApplicationVersion)
			assert.Equal(t, applicationProfile, td.ApplicationProfile)
			assert.Equal(t, trace.ApplicationType, td.ApplicationType)
			assert.NotNil(t, td.ContextId)
			assert.Equal(t, functionName, td.ContextName)
			assert.Equal(t, trace.ExecutionContext, td.ContextType)

			st, err1 := time.Parse(trace.TimeFormat, td.StartTime)
			et, err2 := time.Parse(trace.TimeFormat, td.EndTime)
			if err1 != nil || err2 != nil {
				fmt.Println("err1: ", err1)
				fmt.Println("err2: ", err2)
			}
			assert.True(t, invocationStartTime.Before(st) || invocationStartTime.Equal(st))
			assert.True(t, st.Before(et))
			assert.True(t, et.Before(invocationEndTime) || et.Equal(invocationEndTime))
			assert.True(t, int64(duration) <= td.Duration)

			//Trace Audit Info
			ai := td.AuditInfo
			assert.Equal(t, functionName, ai[trace.AuditInfoContextName])
			assert.NotNil(t, ai[trace.AuditInfoId])
			assert.Equal(t, td.StartTime, ai[trace.AuditInfoOpenTime])
			assert.Equal(t, td.EndTime, ai[trace.AuditInfoCloseTime])

			req := json.RawMessage(testCase.input)

			//Trace Properties
			props := td.Properties
			assert.Equal(t, req, props[trace.AuditInfoPropertiesRequest])
			assert.Equal(t, coldStart, props[trace.AuditInfoPropertiesColdStart])
			assert.Equal(t, defaultRegion, props[trace.AuditInfoPropertiesFunctionRegion])
			assert.Equal(t, memoryLimit, props[trace.AuditInfoPropertiesFunctionMemoryLimit])

			if testCase.expected.err != nil {
				assert.Equal(t, testCase.expected.err, err)

				assert.Equal(t, 1, len(td.Errors))
				assert.Equal(t, errorType, td.ThrownError)
				assert.Equal(t, generatedError, td.ThrownErrorMessage)

				assert.Equal(t, 1, len((ai[trace.AuditInfoErrors]).([]interface{})))
				assert.Nil(t, props[trace.AuditInfoPropertiesResponse])

				errorInfo := ai[trace.AuditInfoThrownError].(trace.ErrorInfo)
				assert.Equal(t, generatedError, errorInfo.ErrMessage)
				assert.Equal(t, errorType, errorInfo.ErrType)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, testCase.expected.val, response)

				assert.Equal(t, testCase.expected.val, props[trace.AuditInfoPropertiesResponse])
				assert.Nil(t, td.Errors)
				assert.Nil(t, td.ThrownError)
				assert.Nil(t, td.ThrownErrorMessage)
				assert.Nil(t, ai[trace.AuditInfoErrors])
				assert.Nil(t, ai[trace.AuditInfoThrownError])
			}

			cleanEnvironment()
			coldStart = "false"
		})
	}
}

func TestPanic(t *testing.T) {
	hello := func(s string) string {
		time.Sleep(time.Millisecond * duration)
		panic(errors.New(generatedPanic))
		return fmt.Sprintf("Happy monitoring with %s!", s)
	}

	testCases := []struct {
		name     string
		input    string
		expected expectedPanic
		handler  interface{}
	}{
		{
			name:     "Panic Test",
			input:    `"Thundra"`,
			expected: expectedPanic{"", nil, errors.New(generatedPanic)},
			handler: func(name string) (string, error) {
				return hello(name), nil
			},
		},
		{
			input:    `"Thundra"`,
			expected: expectedPanic{"Thundra works!", nil,errors.New(generatedPanic)},
			handler: func(name string) (string, error) {
				return hello(name), nil
			},
		},
		{
			input:    `"Thundra"`,
			expected: expectedPanic{"Thundra works!", nil,errors.New(generatedPanic)},
			handler: func(ctx context.Context, name string) (string, error) {
				return hello(name), nil
			},
		},

	}
	for i, testCase := range testCases {
		t.Run(fmt.Sprintf("testCase[%d] %s", i, testCase.name), func(t *testing.T) {
			prepareEnvironment()

			r := new(MockReporter)
			r.On("report").Return()
			r.On("clear").Return()
			r.On("collect", mock.Anything).Return()

			tf := trace.TraceFactory{}
			tr := tf.Create()
			th := NewBuilder().AddPlugin(tr).SetReporter(r).Build()
			lambdaHandler := WrapLambdaHandler(testCase.handler, th)
			invocationStartTime := time.Now().Round(time.Millisecond)

			defer func() {
				if rec := recover(); rec != nil {
					invocationEndTime := time.Now().Round(time.Millisecond)

					//Monitor Data
					msg, ok := r.messageQueue[0].(plugin.Message)
					if !ok {
						fmt.Println("Collector message can't be casted to pluginMessage")
					}
					assert.Equal(t, trace.TraceDataType, msg.Type)
					assert.Equal(t, TestApiKey, msg.ApiKey)
					assert.Equal(t, dataFormatVersion, msg.DataFormatVersion)

					//Trace Data
					td, ok := msg.Data.(trace.TraceData)
					if !ok {
						fmt.Println("Can not convert to trace data")
					}
					assert.NotNil(t, td.Id)
					assert.Equal(t, functionName, td.ApplicationName)
					assert.Equal(t, "", td.ApplicationId)
					assert.Equal(t, functionVersion, td.ApplicationVersion)
					assert.Equal(t, applicationProfile, td.ApplicationProfile)
					assert.Equal(t, trace.ApplicationType, td.ApplicationType)
					assert.NotNil(t, td.ContextId)
					assert.Equal(t, functionName, td.ContextName)
					assert.Equal(t, trace.ExecutionContext, td.ContextType)

					st, err1 := time.Parse(trace.TimeFormat, td.StartTime)
					et, err2 := time.Parse(trace.TimeFormat, td.EndTime)
					if err1 != nil || err2 != nil {
						fmt.Println("err1: ", err1)
						fmt.Println("err2: ", err2)
					}
					assert.True(t, invocationStartTime.Before(st) || invocationStartTime.Equal(st))
					assert.True(t, st.Before(et))
					assert.True(t, et.Before(invocationEndTime) || et.Equal(invocationEndTime))
					assert.True(t, int64(duration) <= td.Duration)

					assert.Equal(t, 1, len(td.Errors))
					assert.Equal(t, errorType, td.ThrownError)
					assert.Equal(t, generatedPanic, td.ThrownErrorMessage)

					//Trace Audit Info
					ai := td.AuditInfo
					assert.Equal(t, functionName, ai[trace.AuditInfoContextName])
					assert.NotNil(t, ai[trace.AuditInfoId])
					assert.Equal(t, td.StartTime, ai[trace.AuditInfoOpenTime])
					assert.Equal(t, td.EndTime, ai[trace.AuditInfoCloseTime])

					panicInfo := ai[trace.AuditInfoThrownError].(trace.PanicInfo)
					assert.Equal(t, generatedPanic, panicInfo.ErrMessage)
					assert.Equal(t, errorType, panicInfo.ErrType)
					assert.NotNil(t, panicInfo.StackTrace)

					req := json.RawMessage(testCase.input)

					//Trace Properties
					props := td.Properties
					assert.Equal(t, req, props[trace.AuditInfoPropertiesRequest])
					assert.Equal(t, coldStart, props[trace.AuditInfoPropertiesColdStart])
					assert.Equal(t, defaultRegion, props[trace.AuditInfoPropertiesFunctionRegion])
					assert.Equal(t, memoryLimit, props[trace.AuditInfoPropertiesFunctionMemoryLimit])

					assert.Equal(t, 1, len((ai[trace.AuditInfoErrors]).([]interface{})))
					assert.Nil(t, props[trace.AuditInfoPropertiesResponse])

					cleanEnvironment()
					coldStart = "false"
				}
			}()
			lambdaHandler.InvokeWithoutSerialization(context.TODO(), []byte(testCase.input))
		})
	}
}

func prepareEnvironment() {
	lambdacontext.FunctionName = functionName
	lambdacontext.MemoryLimitInMB = memoryLimit
	lambdacontext.FunctionVersion = functionVersion
	os.Setenv(trace.ThundraApplicationProfile, applicationProfile)
	os.Setenv(trace.AwsDefaultRegion, defaultRegion)
	SetApiKey(TestApiKey)
}

func cleanEnvironment() {
	lambdacontext.FunctionName = ""
	lambdacontext.MemoryLimitInMB = 0
	lambdacontext.FunctionVersion = ""
	os.Clearenv()
}

//TODO Solve thundra.LambdaFunction type problem
func (handler LambdaFunction) InvokeWithoutSerialization(ctx context.Context, payload []byte) (interface{}, error, time.Time, time.Time) {
	startTime := time.Now().Round(time.Millisecond)
	response, err := handler(ctx, payload)
	endTime := time.Now().Round(time.Millisecond)
	if err != nil {
		return nil, err, startTime, endTime
	}
	return response, nil, startTime, endTime
}
