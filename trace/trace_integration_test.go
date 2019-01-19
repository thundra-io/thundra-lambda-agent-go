package trace

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/thundra-io/thundra-lambda-agent-go/agent"
	"github.com/thundra-io/thundra-lambda-agent-go/plugin"
	"github.com/thundra-io/thundra-lambda-agent-go/test"
)

const (
	duration     = 500
	testAPIKey   = "testApiKey"
	errorMessage = "Error Message"
	errorKind    = "errorString"
	panicMessage = "Panic Message"
)

var coldStart = true

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
			expected: expected{"", errors.New(errorMessage)},
			handler: func() error {
				time.Sleep(time.Millisecond * duration)
				return errors.New(errorMessage)
			},
		},
		{
			input:    `"Thundra"`,
			expected: expected{"", errors.New(errorMessage)},
			handler: func() (interface{}, error) {
				time.Sleep(time.Millisecond * duration)
				return nil, errors.New(errorMessage)
			},
		},
		{
			input:    `"Thundra"`,
			expected: expected{"", errors.New(errorMessage)},
			handler: func(e interface{}) (interface{}, error) {
				time.Sleep(time.Millisecond * duration)
				return nil, errors.New(errorMessage)
			},
		},
		{
			input:    `"Thundra"`,
			expected: expected{"", errors.New(errorMessage)},
			handler: func(ctx context.Context, e interface{}) (interface{}, error) {
				time.Sleep(time.Millisecond * duration)
				return nil, errors.New(errorMessage)
			},
		},
	}
	for i, testCase := range testCases {
		t.Run(fmt.Sprintf("testCase[%d] %s", i, testCase.name), func(t *testing.T) {
			test.PrepareEnvironment()

			r := test.NewMockReporter()
			tr := New()
			a := agent.New().AddPlugin(tr).SetReporter(r)
			lambdaHandler := a.Wrap(testCase.handler)
			h := lambdaHandler.(func(context.Context, json.RawMessage) (interface{}, error))
			f := lambdaFunction(h)
			invocationStartTime := plugin.GetTimestamp()
			response, errVal := f(context.TODO(), []byte(testCase.input))
			invocationEndTime := plugin.GetTimestamp()

			//Monitor Data
			msg, err := getWrappedTraceData(r.MessageQueue)
			if err != nil {
				fmt.Println(err)
				return
			}
			assert.Equal(t, traceType, msg.Type)
			assert.Equal(t, plugin.DataModelVersion, msg.DataModelVersion)

			//Trace Data
			td, ok := msg.Data.(traceDataModel)
			if !ok {
				fmt.Println("Can not convert to trace data")
			}
			assert.NotNil(t, td.ID)
			assert.Equal(t, traceType, td.Type)
			assert.Equal(t, plugin.AgentVersion, td.AgentVersion)
			assert.Equal(t, plugin.DataModelVersion, td.DataModelVersion)
			assert.Equal(t, test.AppId, td.ApplicationID)
			assert.Equal(t, plugin.ApplicationDomainName, td.ApplicationDomainName)
			assert.Equal(t, plugin.ApplicationClassName, td.ApplicationClassName)
			assert.Equal(t, test.FunctionName, td.ApplicationName)
			assert.Equal(t, test.FunctionVersion, td.ApplicationVersion)
			assert.Equal(t, test.ApplicationStage, td.ApplicationStage)
			assert.Equal(t, plugin.ApplicationRuntime, td.ApplicationRuntime)
			assert.Equal(t, plugin.ApplicationRuntimeVersion, td.ApplicationRuntimeVersion)
			assert.NotNil(t, td.ApplicationTags)

			assert.NotNil(t, td.RootSpanID)

			assert.True(t, invocationStartTime <= td.StartTimestamp)
			assert.True(t, td.StartTimestamp < td.FinishTimestamp)
			assert.True(t, td.FinishTimestamp <= invocationEndTime)
			assert.True(t, int64(duration) <= td.Duration)

			//Tags
			assert.Equal(t, testCase.input, td.Tags[plugin.AwsLambdaInvocationRequest])
			assert.Equal(t, test.LogGroupName, td.Tags[plugin.AwsLambdaLogGroupName])
			assert.Equal(t, test.LogStreamName, td.Tags[plugin.AwsLambdaLogStreamName])
			assert.Equal(t, test.MemoryLimit, td.Tags[plugin.AwsLambdaMemoryLimit])
			assert.Equal(t, test.FunctionName, td.Tags[plugin.AwsLambdaName])
			assert.Equal(t, test.Region, td.Tags[plugin.AwsRegion])
			assert.Equal(t, false, td.Tags[plugin.AwsLambdaInvocationTimeout])
			assert.Equal(t, coldStart, td.Tags[plugin.AwsLambdaInvocationColdStart])

			if testCase.expected.err != nil {
				assert.Equal(t, testCase.expected.err, errVal)
				assert.Equal(t, true, td.Tags[plugin.AwsError])
				assert.Equal(t, errorKind, td.Tags[plugin.AwsErrorKind])
				assert.Equal(t, errorMessage, td.Tags[plugin.AwsErrorMessage])
			} else {
				assert.Equal(t, testCase.expected.val, response)
				assert.Equal(t, testCase.expected.val, td.Tags[plugin.AwsLambdaInvocationResponse])
				assert.NoError(t, errVal)
				assert.Nil(t, td.Tags[plugin.AwsError])
				assert.Nil(t, td.Tags[plugin.AwsErrorKind])
				assert.Nil(t, td.Tags[plugin.AwsErrorMessage])
				assert.Nil(t, td.Tags[plugin.AwsErrorStack])
			}

			test.CleanEnvironment()
			coldStart = false
		})
	}
}

func TestPanic(t *testing.T) {
	hello := func(s string) string {
		time.Sleep(time.Millisecond * duration)
		panic(errors.New(panicMessage))
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
			expected: expectedPanic{"", nil, errors.New(panicMessage)},
			handler: func(name string) (string, error) {
				return hello(name), nil
			},
		},
		{
			input:    `"Thundra"`,
			expected: expectedPanic{"Thundra works!", nil, errors.New(panicMessage)},
			handler: func(name string) (string, error) {
				return hello(name), nil
			},
		},
		{
			input:    `"Thundra"`,
			expected: expectedPanic{"Thundra works!", nil, errors.New(panicMessage)},
			handler: func(ctx context.Context, name string) (string, error) {
				return hello(name), nil
			},
		},
	}
	for i, testCase := range testCases {
		t.Run(fmt.Sprintf("testCase[%d] %s", i, testCase.name), func(t *testing.T) {
			test.PrepareEnvironment()

			r := test.NewMockReporter()
			tr := New()
			a := agent.New().AddPlugin(tr).SetReporter(r)
			lambdaHandler := a.Wrap(testCase.handler)
			invocationStartTime := plugin.GetTimestamp()

			defer func() {
				if rec := recover(); rec != nil {
					invocationEndTime := plugin.GetTimestamp()

					//Monitor Data
					msg, err := getWrappedTraceData(r.MessageQueue)
					if err != nil {
						fmt.Println(err)
						return
					}
					assert.Equal(t, traceType, msg.Type)
					assert.Equal(t, plugin.DataModelVersion, msg.DataModelVersion)

					//Trace Data
					td, ok := msg.Data.(traceDataModel)
					if !ok {
						fmt.Println("Can not convert to trace data")
					}
					assert.NotNil(t, td.ID)
					assert.Equal(t, traceType, td.Type)
					assert.Equal(t, plugin.AgentVersion, td.AgentVersion)
					assert.Equal(t, plugin.DataModelVersion, td.DataModelVersion)
					assert.Equal(t, test.AppId, td.ApplicationID)
					assert.Equal(t, plugin.ApplicationDomainName, td.ApplicationDomainName)
					assert.Equal(t, plugin.ApplicationClassName, td.ApplicationClassName)
					assert.Equal(t, test.FunctionName, td.ApplicationName)
					assert.Equal(t, test.FunctionVersion, td.ApplicationVersion)
					assert.Equal(t, test.ApplicationStage, td.ApplicationStage)
					assert.Equal(t, plugin.ApplicationRuntime, td.ApplicationRuntime)
					assert.Equal(t, plugin.ApplicationRuntimeVersion, td.ApplicationRuntimeVersion)
					assert.NotNil(t, td.ApplicationTags)

					assert.NotNil(t, td.RootSpanID)

					assert.True(t, invocationStartTime <= td.StartTimestamp)
					assert.True(t, td.StartTimestamp < td.FinishTimestamp)
					assert.True(t, td.FinishTimestamp <= invocationEndTime)
					assert.True(t, int64(duration) <= td.Duration)

					//Tags
					assert.Equal(t, testCase.input, td.Tags[plugin.AwsLambdaInvocationRequest])
					assert.Equal(t, test.LogGroupName, td.Tags[plugin.AwsLambdaLogGroupName])
					assert.Equal(t, test.LogStreamName, td.Tags[plugin.AwsLambdaLogStreamName])
					assert.Equal(t, test.MemoryLimit, td.Tags[plugin.AwsLambdaMemoryLimit])
					assert.Equal(t, test.FunctionName, td.Tags[plugin.AwsLambdaName])
					assert.Equal(t, test.Region, td.Tags[plugin.AwsRegion])
					assert.Equal(t, false, td.Tags[plugin.AwsLambdaInvocationTimeout])
					assert.Equal(t, coldStart, td.Tags[plugin.AwsLambdaInvocationColdStart])

					//Panic
					assert.Equal(t, true, td.Tags[plugin.AwsError])
					assert.Equal(t, errorKind, td.Tags[plugin.AwsErrorKind])
					assert.Equal(t, panicMessage, td.Tags[plugin.AwsErrorMessage])
					assert.Nil(t, td.Tags[plugin.AwsLambdaInvocationResponse])

					test.CleanEnvironment()
					coldStart = false
				}
			}()
			h := lambdaHandler.(func(context.Context, json.RawMessage) (interface{}, error))
			f := lambdaFunction(h)
			f(context.TODO(), []byte(testCase.input))
		})
	}
}

func TestTimeout(t *testing.T) {
	const timeoutDuration = 1
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

		r := test.NewMockReporter()
		tr := New()
		a := agent.New().AddPlugin(tr).SetReporter(r)
		lambdaHandler := a.Wrap(testCase[0].handler)
		h := lambdaHandler.(func(context.Context, json.RawMessage) (interface{}, error))
		f := lambdaFunction(h)

		d := time.Now().Add(timeoutDuration * time.Second)
		ctx, cancel := context.WithDeadline(context.TODO(), d)
		defer cancel()

		f(ctx, []byte(testCase[0].input))
		// Code doesn't wait goroutines to finish.
		//Monitor Data
		msg, err := getWrappedTraceData(r.MessageQueue)
		if err != nil {
			fmt.Println(err)
			return
		}

		fmt.Println("Current message types in the stack:")
		for _, m := range r.MessageQueue {
			fmt.Println(m.Type)
		}
		//Trace Data
		td, ok := msg.Data.(traceDataModel)
		if !ok {
			fmt.Println("Can not convert to trace data")
		}

		assert.Equal(t, true, td.Tags[plugin.AwsLambdaInvocationTimeout])
	})

}

func getWrappedTraceData(monitoringDataWrappers []plugin.MonitoringDataWrapper) (*plugin.MonitoringDataWrapper, error) {
	for _, m := range monitoringDataWrappers {
		if m.Type == traceType {
			return &m, nil
		}
	}
	return nil, errors.New("trace Data Wrapper is not found")
}

type lambdaFunction func(context.Context, json.RawMessage) (interface{}, error)
