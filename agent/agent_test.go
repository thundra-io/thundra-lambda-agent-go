package agent

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"sync"
	"testing"

	"github.com/stretchr/testify/mock"
	"github.com/thundra-io/thundra-lambda-agent-go/plugin"
	"github.com/thundra-io/thundra-lambda-agent-go/test"
)

const (
	generatedError = "Generated Error"
	testDataType   = "TestDataType"
)

type MockPlugin struct {
	mock.Mock
}

func (t *MockPlugin) IsEnabled() bool {
	return true
}
func (t *MockPlugin) BeforeExecution(ctx context.Context, request json.RawMessage, wg *sync.WaitGroup) {
	defer wg.Done()
	t.Called(ctx, request, wg)
}
func (t *MockPlugin) AfterExecution(ctx context.Context, request json.RawMessage, response interface{}, err interface{}) []plugin.MonitoringDataWrapper {
	t.Called(ctx, request, response, err)
	return []plugin.MonitoringDataWrapper{}
}
func (t *MockPlugin) OnPanic(ctx context.Context, request json.RawMessage, err interface{}, stackTrace []byte) []plugin.MonitoringDataWrapper {
	t.Called(ctx, request, err, stackTrace)
	return []plugin.MonitoringDataWrapper{}
}

func TestExecutePreHooks(t *testing.T) {
	mT := new(MockPlugin)
	th := New().AddPlugin(mT)

	ctx := context.TODO()
	req := createRawMessage()

	mT.On("BeforeExecution", ctx, req, mock.Anything, mock.Anything).Return()
	th.ExecutePreHooks(ctx, req)
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
	var err2 = errors.New("Error")

	r := test.NewMockReporter()

	mT := new(MockPlugin)
	mT.On("AfterExecution", ctx, req, resp, err1, mock.Anything).Return()

	th := New().AddPlugin(mT).SetReporter(r)
	th.ExecutePostHooks(ctx, req, resp, err1)
	th.ExecutePostHooks(ctx, req, resp, err2)

	mT.AssertExpectations(t)

	// Should only be called once because it is already reported
	mT.AssertNumberOfCalls(t, "AfterExecution", 1)
	r.AssertExpectations(t)
}