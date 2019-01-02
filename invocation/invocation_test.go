package invocation

import (
	"context"
	"errors"
	"fmt"
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

func TestInvocationData_BeforeExecution(t *testing.T) {
	i := New()
	prevTime := plugin.GetTimestamp()

	wg := sync.WaitGroup{}
	wg.Add(1)
	i.BeforeExecution(context.TODO(), nil, &wg)

	assert.True(t, prevTime <= i.span.startTimestamp)
}

func TestInvocationData_AfterExecution(t *testing.T) {
	i := New()
	invocationCount = 0
	prevTime := plugin.GetTimestamp()
	i.span.startTimestamp = prevTime

	wg := sync.WaitGroup{}
	wg.Add(1)
	data := i.AfterExecution(context.TODO(), nil, nil, nil)
	d,ok := data[0].Data.(invocationData)
	if !ok{
		fmt.Println("Can not convert to invocationData")
	}

	now := plugin.GetTimestamp()
	assert.True(t, prevTime <= d.FinishTimestamp)
	assert.True(t, d.FinishTimestamp <= now)
	assert.True(t, d.Duration <= now-prevTime)
	assert.True(t, d.ColdStart)
	assert.False(t, d.Timeout)
	assert.Equal(t, invocationType, data[0].Type)
}

func TestInvocationData_AfterExecutionWithError(t *testing.T) {
	const testErrorMessage = "test Error"
	const testErrorType = "errorString"
	i := New()
	err := errors.New(testErrorMessage)

	wg := sync.WaitGroup{}
	wg.Add(1)
	data := i.AfterExecution(context.TODO(), nil, nil, err)
	d,ok := data[0].Data.(invocationData)
	if !ok{
		fmt.Println("Can not convert to invocationData")
	}

	assert.True(t, d.Erroneous)
	assert.Equal(t, testErrorMessage, d.ErrorMessage)
	assert.Equal(t, testErrorType, d.ErrorType)
}

func TestInvocationData_OnPanic(t *testing.T) {
	i := New()
	invocationCount = 0
	prevTime := plugin.GetTimestamp()
	i.span.startTimestamp = prevTime
	err := errors.New(testErrorMessage)

	wg := sync.WaitGroup{}
	wg.Add(1)
	data := i.OnPanic(context.TODO(), nil, err, nil)
	d,ok := data[0].Data.(invocationData)
	if !ok{
		fmt.Println("Can not convert to invocationData")
	}

	now := plugin.GetTimestamp()
	assert.True(t, prevTime <= d.FinishTimestamp)
	assert.True(t, d.FinishTimestamp <= now)
	assert.True(t, d.Duration <= now-prevTime)
	assert.True(t, d.Erroneous)
	assert.Equal(t, testErrorMessage, d.ErrorMessage)
	assert.Equal(t, testErrorType, d.ErrorType)
	assert.True(t, d.ColdStart)
	assert.False(t, d.Timeout)
	assert.Equal(t, invocationType, data[0].Type)
}


func TestPrepareDataStaticFields(t *testing.T) {
	test.PrepareEnvironment()
	i := New()
	data := i.prepareData(context.TODO())
	assert.NotNil(t, data.ID)
	assert.Equal(t, invocationType, data.Type)
	assert.Equal(t, plugin.AgentVersion, data.AgentVersion)
	assert.Equal(t, plugin.DataModelVersion, data.DataModelVersion)
	assert.Equal(t, test.AppId, data.ApplicationID)
	assert.Equal(t, plugin.ApplicationDomainName, data.ApplicationDomainName)
	assert.Equal(t, plugin.ApplicationClassName, data.ApplicationClassName)
	assert.Equal(t, plugin.FunctionName, data.ApplicationName)
	assert.Equal(t, plugin.ApplicationVersion, data.ApplicationVersion)
	assert.Equal(t, plugin.ApplicationStage, data.ApplicationStage)
	assert.Equal(t, plugin.ApplicationRuntime, data.ApplicationRuntime)
	assert.Equal(t, plugin.ApplicationRuntimeVersion, data.ApplicationRuntimeVersion)
	assert.NotNil(t, data.ApplicationTags)

	assert.Equal(t, functionPlatform, data.FunctionPlatform)
	assert.Equal(t, plugin.FunctionName, data.FunctionName)
	assert.Equal(t, plugin.FunctionRegion, data.FunctionRegion)
	assert.NotNil(t, data.Tags)

	test.CleanEnvironment()
}