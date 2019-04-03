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
	"github.com/thundra-io/thundra-lambda-agent-go/application"
	"github.com/thundra-io/thundra-lambda-agent-go/constants"
	"github.com/thundra-io/thundra-lambda-agent-go/plugin"
	"github.com/thundra-io/thundra-lambda-agent-go/test"
	"github.com/thundra-io/thundra-lambda-agent-go/utils"
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
			invocationStartTime := utils.GetTimestamp()
			response, errVal := f(context.TODO(), []byte(testCase.input))
			invocationEndTime := utils.GetTimestamp()

			//Monitor Data
			msg, err := getWrappedTraceData(r.MessageQueue)
			if err != nil {
				fmt.Println(err)
				return
			}
			assert.Equal(t, traceType, msg.Type)
			assert.Equal(t, constants.DataModelVersion, msg.DataModelVersion)

			//Trace Data
			td, ok := msg.Data.(traceDataModel)
			if !ok {
				fmt.Println("Can not convert to trace data")
			}
			assert.NotNil(t, td.ID)
			assert.Equal(t, traceType, td.Type)
			assert.Equal(t, constants.AgentVersion, td.AgentVersion)
			assert.Equal(t, constants.DataModelVersion, td.DataModelVersion)
			assert.Equal(t, test.AppID, td.ApplicationID)
			assert.Equal(t, application.ApplicationDomainName, td.ApplicationDomainName)
			assert.Equal(t, application.ApplicationClassName, td.ApplicationClassName)
			assert.Equal(t, test.ApplicationName, td.ApplicationName)
			assert.Equal(t, test.FunctionVersion, td.ApplicationVersion)
			assert.Equal(t, test.ApplicationStage, td.ApplicationStage)
			assert.Equal(t, application.ApplicationRuntime, td.ApplicationRuntime)
			assert.Equal(t, application.ApplicationRuntimeVersion, td.ApplicationRuntimeVersion)
			assert.NotNil(t, td.ApplicationTags)

			assert.NotNil(t, td.RootSpanID)

			assert.True(t, invocationStartTime <= td.StartTimestamp)
			assert.True(t, td.StartTimestamp < td.FinishTimestamp)
			assert.True(t, td.FinishTimestamp <= invocationEndTime)
			assert.True(t, int64(duration) <= td.Duration)


			if testCase.expected.err != nil {
				assert.Equal(t, testCase.expected.err, errVal)
			} else {
				assert.Equal(t, testCase.expected.val, response)
				assert.NoError(t, errVal)
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
			invocationStartTime := utils.GetTimestamp()

			defer func() {
				if rec := recover(); rec != nil {
					invocationEndTime := utils.GetTimestamp()

					//Monitor Data
					msg, err := getWrappedTraceData(r.MessageQueue)
					if err != nil {
						fmt.Println(err)
						return
					}
					assert.Equal(t, traceType, msg.Type)
					assert.Equal(t, constants.DataModelVersion, msg.DataModelVersion)

					//Trace Data
					td, ok := msg.Data.(traceDataModel)
					if !ok {
						fmt.Println("Can not convert to trace data")
					}
					assert.NotNil(t, td.ID)
					assert.Equal(t, traceType, td.Type)
					assert.Equal(t, constants.AgentVersion, td.AgentVersion)
					assert.Equal(t, constants.DataModelVersion, td.DataModelVersion)
					assert.Equal(t, test.AppID, td.ApplicationID)
					assert.Equal(t, application.ApplicationDomainName, td.ApplicationDomainName)
					assert.Equal(t, application.ApplicationClassName, td.ApplicationClassName)
					assert.Equal(t, test.ApplicationName, td.ApplicationName)
					assert.Equal(t, test.FunctionVersion, td.ApplicationVersion)
					assert.Equal(t, test.ApplicationStage, td.ApplicationStage)
					assert.Equal(t, application.ApplicationRuntime, td.ApplicationRuntime)
					assert.Equal(t, application.ApplicationRuntimeVersion, td.ApplicationRuntimeVersion)
					assert.NotNil(t, td.ApplicationTags)

					assert.NotNil(t, td.RootSpanID)

					assert.True(t, invocationStartTime <= td.StartTimestamp)
					assert.True(t, td.StartTimestamp < td.FinishTimestamp)
					assert.True(t, td.FinishTimestamp <= invocationEndTime)
					assert.True(t, int64(duration) <= td.Duration)

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
		time.Sleep(time.Second * 2 * time.Duration(timeoutDuration))
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

		msg, err := getRootSpanData(r.MessageQueue)
		if err != nil {
			fmt.Println(err)
			return
		}

		//Trace Data
		rsd, ok := msg.Data.(spanDataModel)
		if !ok {
			fmt.Println("Can not convert to trace data")
		}

		assert.Equal(t, true, rsd.Tags[constants.AwsLambdaInvocationTimeout])
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

func getRootSpanData(monitoringDataWrappers []plugin.MonitoringDataWrapper) (*plugin.MonitoringDataWrapper, error) {
	for _, m := range monitoringDataWrappers {
		if m.Type == spanType {
			return &m, nil
		}
	}
	return nil, errors.New("Span Data Wrapper is not found")
}

type lambdaFunction func(context.Context, json.RawMessage) (interface{}, error)
