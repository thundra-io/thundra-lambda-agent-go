package trace

import (
	"testing"
	"context"
	"fmt"
	"os"
	"time"
	"errors"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/aws/aws-lambda-go/lambdacontext"

	"github.com/thundra-io/thundra-lambda-agent-go/plugin"
	"github.com/thundra-io/thundra-lambda-agent-go/test"
	"github.com/thundra-io/thundra-lambda-agent-go/thundra"
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

type expected struct {
	val string
	err error
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

			r := new(test.MockReporter)
			r.On("Report", TestApiKey).Return()
			r.On("Clear").Return()
			r.On("Collect", mock.Anything).Return()

			tr := &Trace{}
			th := thundra.NewBuilder().AddPlugin(tr).SetReporter(r).SetAPIKey(TestApiKey).Build()
			lambdaHandler := thundra.Wrap(testCase.handler, th)

			invocationStartTime := plugin.MakeTimestamp()
			response, err := lambdaHandler(context.TODO(), []byte(testCase.input))
			invocationEndTime := plugin.MakeTimestamp()

			//Monitor Data
			msg, ok := r.MessageQueue[0].(plugin.Message)
			if !ok {
				fmt.Println("Collector message can't be casted to pluginMessage")
			}
			assert.Equal(t, TraceDataType, msg.Type)
			assert.Equal(t, TestApiKey, msg.ApiKey)
			assert.Equal(t, thundra.DataFormatVersion, msg.DataFormatVersion)

			//Trace Data
			td, ok := msg.Data.(traceData)
			if !ok {
				fmt.Println("Can not convert to trace data")
			}
			assert.NotNil(t, td.Id)
			assert.Equal(t, functionName, td.ApplicationName)
			assert.Equal(t, "", td.ApplicationId)
			assert.Equal(t, functionVersion, td.ApplicationVersion)
			assert.Equal(t, applicationProfile, td.ApplicationProfile)
			assert.Equal(t, plugin.ApplicationType, td.ApplicationType)
			assert.NotNil(t, td.ContextId)
			assert.Equal(t, functionName, td.ContextName)
			assert.Equal(t, executionContext, td.ContextType)

			assert.True(t, invocationStartTime <= td.StartTimestamp)
			assert.True(t, td.StartTimestamp < td.EndTimestamp)
			assert.True(t, td.EndTimestamp <= invocationEndTime)
			assert.True(t, int64(duration) <= td.Duration)

			//Trace Audit Info
			ai := td.AuditInfo
			assert.Equal(t, functionName, ai[auditInfoContextName])
			assert.NotNil(t, ai[auditInfoId])
			assert.Equal(t, td.StartTimestamp, ai[auditInfoOpenTimestamp])
			assert.Equal(t, td.EndTimestamp, ai[auditInfoCloseTimestamp])

			//Trace Properties
			props := td.Properties
			assert.Equal(t, testCase.input, props[auditInfoPropertiesRequest])
			assert.Equal(t, coldStart, props[auditInfoPropertiesColdStart])
			assert.Equal(t, defaultRegion, props[auditInfoPropertiesFunctionRegion])
			assert.Equal(t, memoryLimit, props[auditInfoPropertiesFunctionMemoryLimit])

			if testCase.expected.err != nil {
				assert.Equal(t, testCase.expected.err, err)

				assert.Equal(t, 1, len(td.Errors))
				assert.Equal(t, errorType, td.ThrownError)
				assert.Equal(t, generatedError, td.ThrownErrorMessage)

				assert.Equal(t, 1, len((ai[auditInfoErrors]).([]interface{})))
				assert.Nil(t, props[auditInfoPropertiesResponse])

				errorInfo := ai[auditInfoThrownError].(errorInfo)
				assert.Equal(t, generatedError, errorInfo.ErrMessage)
				assert.Equal(t, errorType, errorInfo.ErrType)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, testCase.expected.val, response)

				assert.Equal(t, testCase.expected.val, props[auditInfoPropertiesResponse])
				assert.Nil(t, td.Errors)
				assert.Nil(t, td.ThrownError)
				assert.Nil(t, td.ThrownErrorMessage)
				assert.Nil(t, ai[auditInfoErrors])
				assert.Nil(t, ai[auditInfoThrownError])
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
			expected: expectedPanic{"Thundra works!", nil, errors.New(generatedPanic)},
			handler: func(name string) (string, error) {
				return hello(name), nil
			},
		},
		{
			input:    `"Thundra"`,
			expected: expectedPanic{"Thundra works!", nil, errors.New(generatedPanic)},
			handler: func(ctx context.Context, name string) (string, error) {
				return hello(name), nil
			},
		},
	}
	for i, testCase := range testCases {
		t.Run(fmt.Sprintf("testCase[%d] %s", i, testCase.name), func(t *testing.T) {
			prepareEnvironment()

			r := new(test.MockReporter)
			r.On("Report", TestApiKey).Return()
			r.On("Clear").Return()
			r.On("Collect", mock.Anything).Return()

			tr := &Trace{}
			th := thundra.NewBuilder().AddPlugin(tr).SetReporter(r).SetAPIKey(TestApiKey).Build()
			lambdaHandler := thundra.Wrap(testCase.handler, th)
			invocationStartTime := plugin.MakeTimestamp()

			defer func() {
				if rec := recover(); rec != nil {
					invocationEndTime := plugin.MakeTimestamp()

					//Monitor Data
					msg, ok := r.MessageQueue[0].(plugin.Message)
					if !ok {
						fmt.Println("Collector message can't be casted to pluginMessage")
					}
					assert.Equal(t, TraceDataType, msg.Type)
					assert.Equal(t, TestApiKey, msg.ApiKey)
					assert.Equal(t, thundra.DataFormatVersion, msg.DataFormatVersion)

					//Trace Data
					td, ok := msg.Data.(traceData)
					if !ok {
						fmt.Println("Can not convert to trace data")
					}
					assert.NotNil(t, td.Id)
					assert.Equal(t, functionName, td.ApplicationName)
					assert.Equal(t, "", td.ApplicationId)
					assert.Equal(t, functionVersion, td.ApplicationVersion)
					assert.Equal(t, applicationProfile, td.ApplicationProfile)
					assert.Equal(t, plugin.ApplicationType, td.ApplicationType)
					assert.NotNil(t, td.ContextId)
					assert.Equal(t, functionName, td.ContextName)
					assert.Equal(t, executionContext, td.ContextType)

					assert.True(t, invocationStartTime <= td.StartTimestamp)
					assert.True(t, td.StartTimestamp < td.EndTimestamp)
					assert.True(t, td.EndTimestamp <= invocationEndTime)
					assert.True(t, int64(duration) <= td.Duration)

					assert.Equal(t, 1, len(td.Errors))
					assert.Equal(t, errorType, td.ThrownError)
					assert.Equal(t, generatedPanic, td.ThrownErrorMessage)

					//Trace Audit Info
					ai := td.AuditInfo
					assert.Equal(t, functionName, ai[auditInfoContextName])
					assert.NotNil(t, ai[auditInfoId])
					assert.Equal(t, td.StartTimestamp, ai[auditInfoOpenTimestamp])
					assert.Equal(t, td.EndTimestamp, ai[auditInfoCloseTimestamp])

					panicInfo := ai[auditInfoThrownError].(panicInfo)
					assert.Equal(t, generatedPanic, panicInfo.ErrMessage)
					assert.Equal(t, errorType, panicInfo.ErrType)
					assert.NotNil(t, panicInfo.StackTrace)

					//Trace Properties
					props := td.Properties
					assert.Equal(t, testCase.input, props[auditInfoPropertiesRequest])
					assert.Equal(t, coldStart, props[auditInfoPropertiesColdStart])
					assert.Equal(t, defaultRegion, props[auditInfoPropertiesFunctionRegion])
					assert.Equal(t, memoryLimit, props[auditInfoPropertiesFunctionMemoryLimit])

					assert.Equal(t, 1, len((ai[auditInfoErrors]).([]interface{})))
					assert.Nil(t, props[auditInfoPropertiesResponse])

					cleanEnvironment()
					coldStart = "false"
				}
			}()
			lambdaHandler(context.TODO(), []byte(testCase.input))
		})
	}
}

func prepareEnvironment() {
	lambdacontext.FunctionName = functionName
	lambdacontext.MemoryLimitInMB = memoryLimit
	lambdacontext.FunctionVersion = functionVersion
	os.Setenv(plugin.ThundraApplicationProfile, applicationProfile)
	os.Setenv(awsDefaultRegion, defaultRegion)
}

func cleanEnvironment() {
	lambdacontext.FunctionName = ""
	lambdacontext.MemoryLimitInMB = 0
	lambdacontext.FunctionVersion = ""
	os.Clearenv()
}
