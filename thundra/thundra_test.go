package thundra

import (
	"testing"
	"fmt"
	"context"
	"errors"
	"encoding/json"
	"sync"
	"runtime/debug"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/thundra-io/thundra-lambda-agent-go/test"
)

const (
	generatedError = "Generated Error"
	testApiKey     = "TestApiKey"
	testDataType   = "TestDataType"
)

func TestWrapper(t *testing.T) {
	hello := func(s string) string {
		return fmt.Sprintf("%s works!", s)
	}
	hellop := func(s *string) *string {
		v := hello(*s)
		return &v
	}

	testCases := []struct {
		name     string
		input    string
		expected expected
		handler  interface{}
	}{
		{
			input:    `"Thundra"`,
			expected: expected{`"Thundra works!"`, nil},
			handler: func(name string) (string, error) {
				return hello(name), nil
			},
		},
		{
			input:    `"Thundra"`,
			expected: expected{`"Thundra works!"`, nil},
			handler: func(name string) (string, error) {
				return hello(name), nil
			},
		},
		{
			input:    `"Thundra"`,
			expected: expected{`"Thundra works!"`, nil},
			handler: func(ctx context.Context, name string) (string, error) {
				return hello(name), nil
			},
		},
		{
			input:    `"Thundra"`,
			expected: expected{`"Thundra works!"`, nil},
			handler: func(name *string) (*string, error) {
				return hellop(name), nil
			},
		},
		{
			input:    `"Thundra"`,
			expected: expected{`"Thundra works!"`, nil},
			handler: func(name *string) (*string, error) {
				return hellop(name), nil
			},
		},
		{
			input:    `"Thundra"`,
			expected: expected{`"Thundra works!"`, nil},
			handler: func(ctx context.Context, name *string) (*string, error) {
				return hellop(name), nil
			},
		},
		{
			input:    `"Thundra"`,
			expected: expected{"", errors.New(generatedError)},
			handler: func() error {
				return errors.New(generatedError)
			},
		},
		{
			input:    `"Thundra"`,
			expected: expected{"", errors.New(generatedError)},
			handler: func() (interface{}, error) {
				return nil, errors.New(generatedError)
			},
		},
		{
			input:    `"Thundra"`,
			expected: expected{"", errors.New(generatedError)},
			handler: func(e interface{}) (interface{}, error) {
				return nil, errors.New(generatedError)
			},
		},
		{
			input:    `"Thundra"`,
			expected: expected{"", errors.New(generatedError)},
			handler: func(ctx context.Context, e interface{}) (interface{}, error) {
				return nil, errors.New(generatedError)
			},
		},
		{
			name:     "basic input struct serialization",
			input:    `{"custom":9001}`,
			expected: expected{`9001`, nil},
			handler: func(event struct{ Custom int }) (int, error) {
				return event.Custom, nil
			},
		},
		{
			name:     "basic output struct serialization",
			input:    `9001`,
			expected: expected{`{"Number":9001}`, nil},
			handler: func(event int) (struct{ Number int }, error) {
				return struct{ Number int }{event}, nil
			},
		},
	}
	for i, testCase := range testCases {
		t.Run(fmt.Sprintf("testCase[%d] %s", i, testCase.name), func(t *testing.T) {
			r := &test.MockReporter{}
			r.On("Report", testApiKey).Return()
			r.On("Clear").Return()
			r.On("Collect", mock.Anything).Return()

			th := NewBuilder().SetReporter(r).SetAPIKey(testApiKey).Build()
			lambdaHandler := Wrap(testCase.handler, th)
			response, err := lambdaHandler.invoke(context.TODO(), []byte(testCase.input))

			if testCase.expected.err != nil {
				assert.Equal(t, testCase.expected.err, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, testCase.expected.val, string(response))
			}
		})
	}
}

func TestInvalidWrappers(t *testing.T) {

	testCases := []struct {
		name     string
		handler  interface{}
		expected error
	}{
		{
			name:     "handler is nil",
			expected: errors.New("handler is nil"),
			handler:  nil,
		},
		{
			name:     "handler kind struct is not func",
			expected: errors.New("handler kind struct is not func"),
			handler:  struct{}{},
		},
		{
			name:     "handlers may not take more than two arguments",
			expected: errors.New("handlers may not take more than two arguments, but handler takes 3"),
			handler: func(n context.Context, x string, y string) error {
				return nil
			},
		},
		{
			name:     "two argument handler does not context as first argument",
			expected: errors.New("handler takes two arguments, but the first is not Context. got string"),
			handler: func(a string, x context.Context) error {
				return nil
			},
		},
		{
			name:     "handler may not return more than two values",
			expected: errors.New("handler may not return more than two values"),
			handler: func() (error, error, error) {
				return nil, nil, nil
			},
		},
		{
			name:     "Error has to be the second value",
			expected: errors.New("handler returns two values, but the second does not implement error"),
			handler: func() (error, string) {
				return nil, "thundra"
			},
		},
		{
			name:     "handler returning a single value does not implement error",
			expected: errors.New("handler returns a single value, but it does not implement error"),
			handler: func() string {
				return "thundra"
			},
		},
		{
			name:     "no return value should not result in error",
			expected: nil,
			handler: func() {
			},
		},
	}
	for i, testCase := range testCases {
		t.Run(fmt.Sprintf("testCase[%d] %s", i, testCase.name), func(t *testing.T) {
			r := &test.MockReporter{}
			r.On("Report", testApiKey).Return()
			r.On("Clear").Return()
			r.On("Collect", mock.Anything).Return()

			th := NewBuilder().SetReporter(r).SetAPIKey(testApiKey).Build()
			lambdaHandler := Wrap(testCase.handler, th)
			_, err := lambdaHandler.invoke(context.TODO(), make([]byte, 0))
			assert.Equal(t, testCase.expected, err)
		})
	}
}

type MockPlugin struct {
	mock.Mock
}

func (t *MockPlugin) BeforeExecution(ctx context.Context, request json.RawMessage, wg *sync.WaitGroup) {
	defer wg.Done()
	t.Called(ctx, request, wg)
}
func (t *MockPlugin) AfterExecution(ctx context.Context, request json.RawMessage, response interface{}, err interface{}) ([]interface{}, string) {
	t.Called(ctx, request, response, err)
	return []interface{}{}, testDataType
}
func (t *MockPlugin) OnPanic(ctx context.Context, request json.RawMessage, err interface{}, stackTrace []byte) ([]interface{}, string) {
	t.Called(ctx, request, err, stackTrace)
	return []interface{}{}, testDataType
}

func TestExecutePreHooks(t *testing.T) {
	mT := new(MockPlugin)
	th := NewBuilder().AddPlugin(mT).Build()

	ctx := context.TODO()
	req := createRawMessage()

	mT.On("BeforeExecution", ctx, req, mock.Anything).Return()
	th.executePreHooks(ctx, req)
	mT.AssertExpectations(t)
}

func createRawMessage() json.RawMessage {
	var req json.RawMessage
	event := struct {
		name string
	}{
		"gandalf",
	}

	req, err := json.Marshal(event)
	if err != nil {
		fmt.Println(err)
	}
	return req
}

func TestExecutePostHooks(t *testing.T) {
	type response struct {
		msg string
	}
	ctx := context.TODO()
	req := createRawMessage()
	resp := response{"Thundra"}
	var err1 error = nil
	var err2 error = errors.New("Error")

	r := new(test.MockReporter)
	mT := new(MockPlugin)
	th := NewBuilder().AddPlugin(mT).SetReporter(r).SetAPIKey(testApiKey).Build()

	mT.On("AfterExecution", ctx, req, resp, err1, mock.Anything).Return()
	mT.On("AfterExecution", ctx, req, resp, err2, mock.Anything).Return()
	r.On("Report", testApiKey).Return()
	r.On("Clear").Return()
	r.On("Collect", mock.Anything).Return()

	th.executePostHooks(ctx, req, resp, err1)
	th.executePostHooks(ctx, req, resp, err2)
	mT.AssertExpectations(t)
	r.AssertExpectations(t)
}

func TestOnPanic(t *testing.T) {
	ctx := context.TODO()
	req := createRawMessage()
	err := errors.New("Generated Error")
	stackTrace := debug.Stack()

	r := new(test.MockReporter)
	mT := new(MockPlugin)
	th := NewBuilder().AddPlugin(mT).SetReporter(r).SetAPIKey(testApiKey).Build()

	mT.On("OnPanic", ctx, req, err, stackTrace, mock.Anything).Return()
	r.On("Report", testApiKey).Return()
	r.On("Clear").Return()
	r.On("Collect", mock.Anything).Return()

	th.onPanic(ctx, req, err, stackTrace)
	mT.AssertExpectations(t)
	r.AssertExpectations(t)
}

func (handler LambdaFunction) invoke(ctx context.Context, payload []byte) ([]byte, error) {
	response, err := handler(ctx, payload)
	if err != nil {
		return nil, err
	}

	responseBytes, err := json.Marshal(response)
	if err != nil {
		return nil, err
	}

	return responseBytes, nil
}

type expected struct {
	val string
	err error
}
