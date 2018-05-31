package invocation

import (
	"context"
	"errors"
	"sync"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/thundra-io/thundra-lambda-agent-go/plugin"
	"github.com/thundra-io/thundra-lambda-agent-go/test"
)

const (
	testErrorMessage = "test Error"
	testErrorType    = "errorString"
)

func TestNewInvocation(t *testing.T) {
	test.PrepareEnvironment()
	i := New()

	assert.NotNil(t, i.Id)
	assert.Equal(t, test.FunctionName, i.ApplicationName)
	assert.Equal(t, test.AppId, i.ApplicationId)
	assert.Equal(t, test.FunctionVersion, i.ApplicationVersion)
	assert.Equal(t, plugin.ApplicationType, i.ApplicationType)
	assert.Equal(t, test.Region, i.Region)
	assert.Equal(t, test.MemoryLimit, i.MemorySize)

	test.CleanEnvironment()
}

func TestInvocationData_BeforeExecution(t *testing.T) {
	i := New()
	prevId := i.Id
	prevTime := plugin.GetTimestamp()

	wg := sync.WaitGroup{}
	wg.Add(1)
	i.BeforeExecution(context.TODO(), nil, &wg)

	assert.NotEqual(t, prevId, i.Id)
	assert.Equal(t, plugin.TransactionId, i.TransactionId)
	assert.True(t, prevTime <= i.StartTimestamp)
}

func TestInvocationData_AfterExecution(t *testing.T) {
	i := New()
	invocationCount = 0
	prevTime := plugin.GetTimestamp()
	i.StartTimestamp = prevTime

	wg := sync.WaitGroup{}
	wg.Add(1)
	_, dataType := i.AfterExecution(context.TODO(), nil, nil, nil)
	now := plugin.GetTimestamp()
	assert.True(t, prevTime <= i.EndTimestamp)
	assert.True(t, i.EndTimestamp <= now)
	assert.True(t, i.Duration <= now-prevTime)
	assert.True(t, i.ColdStart)
	assert.False(t, i.Timeout)
	assert.Equal(t, invocationDataType, dataType)
}

func TestInvocationData_AfterExecutionWithError(t *testing.T) {
	const testErrorMessage = "test Error"
	const testErrorType = "errorString"
	i := New()
	err := errors.New(testErrorMessage)

	wg := sync.WaitGroup{}
	wg.Add(1)
	i.AfterExecution(context.TODO(), nil, nil, err)
	assert.True(t, i.Erroneous)
	assert.Equal(t, testErrorMessage, i.ErrorMessage)
	assert.Equal(t, testErrorType, i.ErrorType)
}

func TestInvocationData_OnPanic(t *testing.T) {
	i := New()
	invocationCount = 0
	prevTime := plugin.GetTimestamp()
	i.StartTimestamp = prevTime
	err := errors.New(testErrorMessage)

	wg := sync.WaitGroup{}
	wg.Add(1)
	_, dataType := i.OnPanic(context.TODO(), nil, err, nil)
	now := plugin.GetTimestamp()
	assert.True(t, prevTime <= i.EndTimestamp)
	assert.True(t, i.EndTimestamp <= now)
	assert.True(t, i.Duration <= now-prevTime)
	assert.True(t, i.Erroneous)
	assert.Equal(t, testErrorMessage, i.ErrorMessage)
	assert.Equal(t, testErrorType, i.ErrorType)
	assert.True(t, i.ColdStart)
	assert.False(t, i.Timeout)
	assert.Equal(t, invocationDataType, dataType)
}
