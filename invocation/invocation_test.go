package invocation

import (
	"context"
	"errors"
	"os"
	"sync"
	"testing"

	"github.com/aws/aws-lambda-go/lambdacontext"
	"github.com/stretchr/testify/assert"
	"github.com/thundra-io/thundra-lambda-agent-go/plugin"
)

const (
	functionName       = "TestFunctionName"
	memoryLimit        = 512
	functionVersion    = "$Version"
	appId              = "$AppId"
	applicationProfile = "TestProfile"
	defaultRegion      = "TestRegion"
	testErrorMessage   = "test Error"
	testErrorType      = "errorString"
)

func TestNewInvocation(t *testing.T) {
	prepareEnvironment()
	i := NewInvocation()

	assert.NotNil(t, i.Id)
	assert.Equal(t, functionName, i.ApplicationName)
	assert.Equal(t, appId, i.ApplicationId)
	assert.Equal(t, functionVersion, i.ApplicationVersion)
	assert.Equal(t, plugin.ApplicationType, i.ApplicationType)
	assert.Equal(t, defaultRegion, i.Region)
	assert.Equal(t, memoryLimit, i.MemorySize)

	cleanEnvironment()
}

func TestInvocationData_BeforeExecution(t *testing.T) {
	i := NewInvocation()
	prevId := i.Id
	transId := plugin.GenerateNewId()
	prevTime := plugin.GetTimestamp()

	wg := sync.WaitGroup{}
	wg.Add(1)
	i.BeforeExecution(context.TODO(), nil, transId, &wg)

	assert.NotEqual(t, prevId, i.Id)
	assert.Equal(t, transId, i.TransactionId)
	assert.True(t, prevTime <= i.StartTimestamp)
}

func TestInvocationData_AfterExecution(t *testing.T) {
	i := NewInvocation()
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
	i := NewInvocation()
	err := errors.New(testErrorMessage)

	wg := sync.WaitGroup{}
	wg.Add(1)
	i.AfterExecution(context.TODO(), nil, nil, err)
	assert.True(t, i.Erroneous)
	assert.Equal(t, testErrorMessage, i.ErrorMessage)
	assert.Equal(t, testErrorType, i.ErrorType)
}

func TestInvocationData_OnPanic(t *testing.T) {
	i := NewInvocation()
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

func prepareEnvironment() {
	lambdacontext.FunctionName = functionName
	lambdacontext.MemoryLimitInMB = memoryLimit
	lambdacontext.FunctionVersion = functionVersion
	lambdacontext.LogStreamName = "[]" + appId
	os.Setenv(plugin.ThundraApplicationProfile, applicationProfile)
	os.Setenv(plugin.AwsDefaultRegion, defaultRegion)
}

func cleanEnvironment() {
	lambdacontext.FunctionName = ""
	lambdacontext.MemoryLimitInMB = 0
	lambdacontext.FunctionVersion = ""
	os.Clearenv()
}
