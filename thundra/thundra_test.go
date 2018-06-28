package thundra

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"runtime/debug"
	"sync"
	"testing"

	"github.com/stretchr/testify/mock"
	"github.com/thundra-io/thundra-lambda-agent-go/test"
)

const (
	generatedError = "Generated Error"
	testApiKey     = "TestApiKey"
	testDataType   = "TestDataType"
)

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
	th := NewBuilder().AddPlugin(mT).SetAPIKey(testApiKey).Build()

	ctx := context.TODO()
	req := createRawMessage()

	mT.On("BeforeExecution", ctx, req, mock.Anything, mock.Anything).Return()
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
	var err1 error
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
	mP := new(MockPlugin)
	th := NewBuilder().AddPlugin(mP).SetReporter(r).SetAPIKey(testApiKey).Build()

	mP.On("OnPanic", ctx, req, err, stackTrace, mock.Anything).Return()
	r.On("Report", testApiKey).Return()
	r.On("Clear").Return()
	r.On("Collect", mock.Anything).Return()

	th.onPanic(ctx, req, err, stackTrace)
	mP.AssertExpectations(t)
	r.AssertExpectations(t)
}

func (handler lambdaFunction) invoke(ctx context.Context, payload []byte) ([]byte, error) {
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
