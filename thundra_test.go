package thundra

import (
	"testing"
	"fmt"
	"context"
	"errors"
	"github.com/stretchr/testify/assert"
	"encoding/json"
	"github.com/stretchr/testify/mock"
	"sync"
)

// Invoke calls the handler, and serializes the response.
// If the underlying handler returned an error, or an error occurs during serialization, error is returned.
func (handler thundraLambdaHandler) Invoke(ctx context.Context, payload []byte) ([]byte, error) {
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

func TestWrapper(t *testing.T) {
	hello := func(s string) string {
		return fmt.Sprintf("Happy monitoring with %s!", s)
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
			expected: expected{`"Happy monitoring with Thundra!"`, nil},
			handler: func(name string) (string, error) {
				return hello(name), nil
			},
		},
		{
			input:    `"Thundra"`,
			expected: expected{`"Happy monitoring with Thundra!"`, nil},
			handler: func(name string) (string, error) {
				return hello(name), nil
			},
		},
		{
			input:    `"Thundra"`,
			expected: expected{`"Happy monitoring with Thundra!"`, nil},
			handler: func(ctx context.Context, name string) (string, error) {
				return hello(name), nil
			},
		},
		{
			input:    `"Thundra"`,
			expected: expected{`"Happy monitoring with Thundra!"`, nil},
			handler: func(name *string) (*string, error) {
				return hellop(name), nil
			},
		},
		{
			input:    `"Thundra"`,
			expected: expected{`"Happy monitoring with Thundra!"`, nil},
			handler: func(name *string) (*string, error) {
				return hellop(name), nil
			},
		},
		{
			input:    `"Thundra"`,
			expected: expected{`"Happy monitoring with Thundra!"`, nil},
			handler: func(ctx context.Context, name *string) (*string, error) {
				return hellop(name), nil
			},
		},
		{
			input:    `"Thundra"`,
			expected: expected{"", errors.New("thundra is dead baby, thundra is dead")},
			handler: func() error {
				return errors.New("thundra is dead baby, thundra is dead")
			},
		},
		{
			input:    `"Thundra"`,
			expected: expected{"", errors.New("thundra is dead baby, thundra is dead")},
			handler: func() (interface{}, error) {
				return nil, errors.New("thundra is dead baby, thundra is dead")
			},
		},
		{
			input:    `"Thundra"`,
			expected: expected{"", errors.New("thundra is dead baby, thundra is dead")},
			handler: func(e interface{}) (interface{}, error) {
				return nil, errors.New("thundra is dead baby, thundra is dead")
			},
		},
		{
			input:    `"Thundra"`,
			expected: expected{"", errors.New("thundra is dead baby, thundra is dead")},
			handler: func(ctx context.Context, e interface{}) (interface{}, error) {
				return nil, errors.New("thundra is dead baby, thundra is dead")
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
			th := CreateNew([]string{})
			lambdaHandler := WrapLambdaHandler(testCase.handler, th)
			response, err := lambdaHandler.Invoke(context.TODO(), []byte(testCase.input))

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
			th := CreateNew([]string{})
			lambdaHandler := WrapLambdaHandler(testCase.handler, th)
			_, err := lambdaHandler.Invoke(context.TODO(), make([]byte, 0))
			assert.Equal(t, testCase.expected, err)
		})
	}
}

type MockPlugin struct {
	mock.Mock
}

type MockedPluginFactory struct{}

func (t *MockedPluginFactory) Create() Plugin {
	return &MockPlugin{}
}

func (t *MockPlugin) BeforeExecution(ctx context.Context, request interface{}, wg *sync.WaitGroup) {
	defer wg.Done()
	t.Called(ctx, request, wg)
}
func (t *MockPlugin) AfterExecution(ctx context.Context, request interface{}, response interface{}, error interface{}, wg *sync.WaitGroup) {
	defer wg.Done()
	t.Called(ctx, request, response, error, wg)
}
func (t *MockPlugin) OnPanic(ctx context.Context, request json.RawMessage, panic *ThundraPanic, wg *sync.WaitGroup) {
	defer wg.Done()
	t.Called(ctx, request, panic, wg)
}

func TestExecutePreHooks(t *testing.T) {
	th := CreateNew([]string{})
	mT := new(MockPlugin)
	th.addPlugin(mT)

	ctx := *new(context.Context)
	req := json.RawMessage{}
	mT.On("BeforeExecution", ctx, req, mock.Anything).Return()

	th.executePreHooks(ctx, req)

	mT.AssertExpectations(t)
}

type MockCollector struct {
	mock.Mock
	msg Message
}

func (c *MockCollector) collect(msg Message) {
	c.Called(msg)
	c.msg = msg
}

func (c *MockCollector) report() {
	c.Called()
}

func (c *MockCollector) clear() {
	c.Called()
}

func TestExecutePostHooks(t *testing.T) {
	type response struct {
		msg string
	}
	//TODO context.TODO()
	ctx := *new(context.Context)
	req := json.RawMessage{}
	resp := response{"Thundra"}
	var err1 error = nil
	var err2 error = errors.New("Error")

	c := new(MockCollector)
	th := createNewWithCollector([]string{}, c)
	mT := new(MockPlugin)
	th.addPlugin(mT)

	mT.On("AfterExecution", ctx, req, resp, err1, mock.Anything).Return()
	mT.On("AfterExecution", ctx, req, resp, err2, mock.Anything).Return()
	c.On("report").Return()
	c.On("clear").Return()

	th.executePostHooks(ctx, req, resp, err1)
	th.executePostHooks(ctx, req, resp, err2)
	mT.AssertExpectations(t)
	c.AssertExpectations(t)
}

func TestOnPanic(t *testing.T) {
	ctx := *new(context.Context)
	req := json.RawMessage{}
	panic := ThundraPanic{
		"Panic Message",
		"runtime/debug.Stack(0xc420043f60, 0x1, 0x1)/n" +
			"/r/usr/local/go/src/runtime/debug/stack.go:24 +0xa7",
		"String Error",
	}

	c := new(MockCollector)
	th := createNewWithCollector([]string{}, c)
	mT := new(MockPlugin)
	th.addPlugin(mT)

	mT.On("OnPanic", ctx, req, &panic, mock.Anything).Return()
	c.On("report").Return()
	c.On("clear").Return()

	th.onPanic(ctx, req, &panic)
	mT.AssertExpectations(t)
	c.AssertExpectations(t)
}
