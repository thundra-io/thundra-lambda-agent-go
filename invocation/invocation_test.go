package invocation

import (
	"context"
	"errors"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/thundra-io/thundra-lambda-agent-go/application"
	"github.com/thundra-io/thundra-lambda-agent-go/constants"
	"github.com/thundra-io/thundra-lambda-agent-go/test"
	"github.com/thundra-io/thundra-lambda-agent-go/utils"
)

const (
	testErrorMessage = "test Error"
	testErrorType    = "errorString"
)

func TestInvocationData_BeforeExecution(t *testing.T) {
	ip := New()
	prevTime := utils.GetTimestamp()

	ip.BeforeExecution(context.TODO(), nil)

	assert.True(t, prevTime <= ip.data.startTimestamp)
}

func TestInvocationData_AfterExecution(t *testing.T) {
	ip := New()
	invocationCount = 0
	prevTime := utils.GetTimestamp()
	ip.data.startTimestamp = prevTime

	data := ip.AfterExecution(context.TODO(), nil, nil, nil)
	d, ok := data[0].Data.(invocationDataModel)
	if !ok {
		fmt.Println("Can not convert to invocationDataModel")
	}

	now := utils.GetTimestamp()
	assert.True(t, prevTime <= d.FinishTimestamp)
	assert.True(t, d.FinishTimestamp <= now)
	assert.True(t, d.Duration <= now-prevTime)
	assert.True(t, d.ColdStart)
	assert.False(t, d.Timeout)
	assert.Equal(t, "Invocation", data[0].Type)

	ClearTags()
}

func TestInvocationData_AfterExecutionWithError(t *testing.T) {
	const testErrorMessage = "test Error"
	const testErrorType = "errorString"
	ip := New()
	err := errors.New(testErrorMessage)

	data := ip.AfterExecution(context.TODO(), nil, nil, err)
	d, ok := data[0].Data.(invocationDataModel)
	if !ok {
		fmt.Println("Can not convert to invocationDataModel")
	}

	assert.True(t, d.Erroneous)
	assert.Equal(t, testErrorMessage, d.ErrorMessage)
	assert.Equal(t, testErrorType, d.ErrorType)

	ClearTags()
}

func TestPrepareDataStaticFields(t *testing.T) {
	test.PrepareEnvironment()
	i := New()
	data := i.prepareData(context.TODO())
	assert.NotNil(t, data.ID)
	assert.Equal(t, "Invocation", data.Type)
	assert.Equal(t, constants.AgentVersion, data.AgentVersion)
	assert.Equal(t, constants.DataModelVersion, data.DataModelVersion)
	assert.Equal(t, test.AppId, data.ApplicationID)
	assert.Equal(t, application.ApplicationDomainName, data.ApplicationDomainName)
	assert.Equal(t, application.ApplicationClassName, data.ApplicationClassName)
	assert.Equal(t, application.FunctionName, data.ApplicationName)
	assert.Equal(t, application.ApplicationVersion, data.ApplicationVersion)
	assert.Equal(t, application.ApplicationStage, data.ApplicationStage)
	assert.Equal(t, application.ApplicationRuntime, data.ApplicationRuntime)
	assert.Equal(t, application.ApplicationRuntimeVersion, data.ApplicationRuntimeVersion)
	assert.NotNil(t, data.ApplicationTags)

	assert.Equal(t, "AWS Lambda", data.FunctionPlatform)
	assert.Equal(t, application.FunctionName, data.FunctionName)
	assert.Equal(t, application.FunctionRegion, data.FunctionRegion)
	assert.NotNil(t, data.Tags)

	test.CleanEnvironment()
	ClearTags()
}

func TestInvocationTags(t *testing.T) {
	tags := map[string]interface{}{
		"boolKey":   true,
		"intKey":    37,
		"floatKey":  3.14,
		"stringKey": "foobar",
		"dictKey": map[string]string{
			"key1": "val1",
			"key2": "val2",
		},
	}

	for k, v := range tags {
		SetTag(k, v)
	}

	assert.Equal(t, len(invocationTags), len(tags))
	assert.Equal(t, invocationTags["boolKey"], tags["boolKey"])
	assert.Equal(t, invocationTags["intKey"], tags["intKey"])
	assert.Equal(t, invocationTags["floatKey"], tags["floatKey"])
	assert.Equal(t, invocationTags["stringKey"], tags["stringKey"])
	assert.Equal(t, invocationTags["dictKey"], tags["dictKey"])

	ClearTags()
}
