package trace

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/thundra-io/thundra-lambda-agent-go/plugin"
	"github.com/thundra-io/thundra-lambda-agent-go/test"
	"github.com/thundra-io/thundra-lambda-agent-go/thundra"
)

const (
	duration       = 500
	testApiKey     = "testApiKey"
	generatedError = "Generated Error"
	errorType      = "errorString"
	generatedPanic = "Generated Panic"
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
			test.PrepareEnvironment()

			r := test.NewMockReporter(testApiKey)
			tr := New()
			th := thundra.NewBuilder().AddPlugin(tr).SetReporter(r).SetAPIKey(testApiKey).Build()
			lambdaHandler := thundra.Wrap(testCase.handler, th)
			h := lambdaHandler.(func(context.Context, json.RawMessage) (interface{}, error))
			f := lambdaFunction(h)

			invocationStartTime := plugin.GetTimestamp()
			response, err := f(context.TODO(), []byte(testCase.input))
			invocationEndTime := plugin.GetTimestamp()

			//Monitor Data
			msg, ok := r.MessageQueue[1].(plugin.Message)
			if !ok {
				fmt.Println("Collector message can't be casted to pluginMessage")
			}
			assert.Equal(t, traceDataType, msg.Type)
			assert.Equal(t, testApiKey, msg.ApiKey)
			assert.Equal(t, "1.2", msg.DataFormatVersion)

			//Trace Data
			td, ok := msg.Data.(traceData)
			if !ok {
				fmt.Println("Can not convert to trace data")
			}
			assert.NotNil(t, td.Id)
			assert.Equal(t, test.FunctionName, td.ApplicationName)
			assert.Equal(t, test.AppId, td.ApplicationId)
			assert.Equal(t, test.FunctionVersion, td.ApplicationVersion)
			assert.Equal(t, test.ApplicationProfile, td.ApplicationProfile)
			assert.Equal(t, plugin.ApplicationType, td.ApplicationType)
			assert.NotNil(t, td.ContextId)
			assert.Equal(t, test.FunctionName, td.ContextName)
			assert.Equal(t, executionContext, td.ContextType)

			assert.True(t, invocationStartTime <= td.StartTimestamp)
			assert.True(t, td.StartTimestamp < td.EndTimestamp)
			assert.True(t, td.EndTimestamp <= invocationEndTime)
			assert.True(t, int64(duration) <= td.Duration)

			//Trace Audit Info
			ai := td.AuditInfo
			assert.Equal(t, test.FunctionName, ai[auditInfoContextName])
			assert.NotNil(t, ai[auditInfoId])
			assert.Equal(t, td.StartTimestamp, ai[auditInfoOpenTimestamp])
			assert.Equal(t, td.EndTimestamp, ai[auditInfoCloseTimestamp])

			//Trace Properties
			props := td.Properties
			assert.Equal(t, testCase.input, props[auditInfoPropertiesRequest])
			assert.Equal(t, coldStart, props[auditInfoPropertiesColdStart])
			assert.Equal(t, test.Region, props[auditInfoPropertiesFunctionRegion])
			assert.Equal(t, test.MemoryLimit, props[auditInfoPropertiesFunctionMemoryLimit])
			assert.Equal(t, test.LogGroupName, props[auditInfoPropertiesLogGroupName])
			assert.Equal(t, test.LogStreamName, props[auditInfoPropertiesLogStreamName])
			assert.NotNil(t, props[auditInfoPropertiesFunctionARN])
			assert.NotNil(t, props[auditInfoPropertiesRequestId])

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

			test.CleanEnvironment()
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
			test.PrepareEnvironment()

			r := test.NewMockReporter(testApiKey)
			tr := New()
			th := thundra.NewBuilder().AddPlugin(tr).SetReporter(r).SetAPIKey(testApiKey).Build()
			lambdaHandler := thundra.Wrap(testCase.handler, th)
			invocationStartTime := plugin.GetTimestamp()

			defer func() {
				if rec := recover(); rec != nil {
					invocationEndTime := plugin.GetTimestamp()

					//Monitor Data
					msg, ok := r.MessageQueue[1].(plugin.Message)
					if !ok {
						fmt.Println("Collector message can't be casted to pluginMessage")
					}
					assert.Equal(t, traceDataType, msg.Type)
					assert.Equal(t, testApiKey, msg.ApiKey)
					assert.Equal(t, "1.2", msg.DataFormatVersion)

					//Trace Data
					td, ok := msg.Data.(traceData)
					if !ok {
						fmt.Println("Can not convert to trace data")
					}
					assert.NotNil(t, td.Id)
					assert.Equal(t, test.FunctionName, td.ApplicationName)
					assert.Equal(t, test.AppId, td.ApplicationId)
					assert.Equal(t, test.FunctionVersion, td.ApplicationVersion)
					assert.Equal(t, test.ApplicationProfile, td.ApplicationProfile)
					assert.Equal(t, plugin.ApplicationType, td.ApplicationType)
					assert.NotNil(t, td.ContextId)
					assert.Equal(t, test.FunctionName, td.ContextName)
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
					assert.Equal(t, test.FunctionName, ai[auditInfoContextName])
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
					assert.Equal(t, test.Region, props[auditInfoPropertiesFunctionRegion])
					assert.Equal(t, test.MemoryLimit, props[auditInfoPropertiesFunctionMemoryLimit])
					assert.Equal(t, test.LogGroupName, props[auditInfoPropertiesLogGroupName])
					assert.Equal(t, test.LogStreamName, props[auditInfoPropertiesLogStreamName])
					assert.NotNil(t, props[auditInfoPropertiesFunctionARN])
					assert.NotNil(t, props[auditInfoPropertiesRequestId])

					assert.Equal(t, 1, len((ai[auditInfoErrors]).([]interface{})))
					assert.Nil(t, props[auditInfoPropertiesResponse])

					test.CleanEnvironment()
					coldStart = "false"
				}
			}()
			h := lambdaHandler.(func(context.Context, json.RawMessage) (interface{}, error))
			f := lambdaFunction(h)
			f(context.TODO(), []byte(testCase.input))
		})
	}
}

func TestTimeout(t *testing.T) {
	const timeoutDuration = 1;
	timeOutFunction := func(s string) string {
		// Let it run longer than timeoutDuration
		time.Sleep(time.Second * 2 * timeoutDuration)
		return fmt.Sprintf("Happy monitoring with %s!", s)
	}

	testCase := []struct {
		name     string
		input    string
		expected expectedPanic
		handler  interface{}
	}{
		{
			name:  "Timeout Test",
			input: `"Thundra"`,
			handler: func(name string) (string, error) {
				return timeOutFunction(name), nil
			},
		},
	}
	t.Run(fmt.Sprintf("testCase[%d] %s", 0, testCase[0].name), func(t *testing.T) {
		test.PrepareEnvironment()

		r := test.NewMockReporter(testApiKey)
		tr := New()
		th := thundra.NewBuilder().AddPlugin(tr).SetReporter(r).SetAPIKey(testApiKey).Build()
		lambdaHandler := thundra.Wrap(testCase[0].handler, th)
		h := lambdaHandler.(func(context.Context, json.RawMessage) (interface{}, error))
		f := lambdaFunction(h)

		d := time.Now().Add(timeoutDuration * time.Second)
		ctx, cancel := context.WithDeadline(context.TODO(), d)
		defer cancel()

		f(ctx, []byte(testCase[0].input))
		// Code doesn't wait goroutines to finish.
		//Monitor Data
		msg, ok := r.MessageQueue[1].(plugin.Message)
		if !ok {
			fmt.Println("Collector message can't be casted to pluginMessage")
		}

		//Trace Data
		td, ok := msg.Data.(traceData)
		if !ok {
			fmt.Println("Can not convert to trace data")
		}

		assert.Equal(t, "true", td.Properties[auditInfoPropertiesTimeout])
	})

}

type lambdaFunction func(context.Context, json.RawMessage) (interface{}, error)
