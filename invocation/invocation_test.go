package invocation

import (
	"context"
	"errors"
	"fmt"
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
	ip := New()
	prevTime := plugin.GetTimestamp()

	ip.BeforeExecution(context.TODO(), nil)

	assert.True(t, prevTime <= ip.data.startTimestamp)
}

func TestInvocationData_AfterExecution(t *testing.T) {
	ip := New()
	invocationCount = 0
	prevTime := plugin.GetTimestamp()
	ip.data.startTimestamp = prevTime

	data := ip.AfterExecution(context.TODO(), nil, nil, nil)
	d, ok := data[0].Data.(invocationDataModel)
	if !ok {
		fmt.Println("Can not convert to invocationDataModel")
	}

	now := plugin.GetTimestamp()
	assert.True(t, prevTime <= d.FinishTimestamp)
	assert.True(t, d.FinishTimestamp <= now)
	assert.True(t, d.Duration <= now-prevTime)
	assert.True(t, d.ColdStart)
	assert.False(t, d.Timeout)
	assert.Equal(t, "Invocation", data[0].Type)
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
}

func TestPrepareDataStaticFields(t *testing.T) {
	test.PrepareEnvironment()
	i := New()
	data := i.prepareData(context.TODO())
	assert.NotNil(t, data.ID)
	assert.Equal(t, "Invocation", data.Type)
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

	assert.Equal(t, "AWS Lambda", data.FunctionPlatform)
	assert.Equal(t, plugin.FunctionName, data.FunctionName)
	assert.Equal(t, plugin.FunctionRegion, data.FunctionRegion)
	assert.NotNil(t, data.Tags)

	test.CleanEnvironment()
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
		SetInvocationTag(k, v)
	}

	assert.Equal(t, len(invocationTags), len(tags))
	assert.Equal(t, invocationTags["boolKey"], tags["boolKey"])
	assert.Equal(t, invocationTags["intKey"], tags["intKey"])
	assert.Equal(t, invocationTags["floatKey"], tags["floatKey"])
	assert.Equal(t, invocationTags["stringKey"], tags["stringKey"])
	assert.Equal(t, invocationTags["dictKey"], fmt.Sprint(tags["dictKey"]))
}
