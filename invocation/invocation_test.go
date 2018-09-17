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
	assert.Equal(t, invocationType, i.Type)
	assert.Equal(t, plugin.AgentVersion, i.AgentVersion)
	assert.Equal(t, plugin.DataModelVersion, i.DataModelVersion)
	assert.Equal(t, test.AppId, i.ApplicationId)
	assert.Equal(t, plugin.ApplicationDomainName, i.ApplicationDomainName)
	assert.Equal(t, plugin.ApplicationClassName, i.ApplicationClassName)
	assert.Equal(t, plugin.FunctionName, i.ApplicationName)
	assert.Equal(t, plugin.ApplicationVersion, i.ApplicationVersion)
	assert.Equal(t, plugin.ApplicationStage, i.ApplicationStage)
	assert.Equal(t, plugin.ApplicationRuntime, i.ApplicationRuntime)
	assert.Equal(t, plugin.ApplicationRuntimeVersion, i.ApplicationRuntimeVersion)
	assert.NotNil(t, i.ApplicationTags)

	assert.Equal(t, functionPlatform, i.FunctionPlatform)
	assert.Equal(t, plugin.FunctionName, i.FunctionName)
	assert.Equal(t, plugin.FunctionRegion, i.FunctionRegion)
	assert.NotNil(t, i.Tags)

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
	assert.True(t, prevTime <= i.FinishTimestamp)
	assert.True(t, i.FinishTimestamp <= now)
	assert.True(t, i.Duration <= now-prevTime)
	assert.True(t, i.ColdStart)
	assert.False(t, i.Timeout)
	assert.Equal(t, invocationType, dataType)
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
	assert.True(t, prevTime <= i.FinishTimestamp)
	assert.True(t, i.FinishTimestamp <= now)
	assert.True(t, i.Duration <= now-prevTime)
	assert.True(t, i.Erroneous)
	assert.Equal(t, testErrorMessage, i.ErrorMessage)
	assert.Equal(t, testErrorType, i.ErrorType)
	assert.True(t, i.ColdStart)
	assert.False(t, i.Timeout)
	assert.Equal(t, invocationType, dataType)
}
