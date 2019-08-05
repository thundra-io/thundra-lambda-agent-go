package invocation

import (
	"context"
	"errors"
	"log"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/thundra-io/thundra-lambda-agent-go/application"
	"github.com/thundra-io/thundra-lambda-agent-go/config"
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

	data, _ := ip.AfterExecution(context.TODO(), nil, nil, nil)
	d, ok := data[0].Data.(invocationDataModel)
	if !ok {
		log.Println("Can not convert to invocationDataModel")
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

	data, _ := ip.AfterExecution(context.TODO(), nil, nil, err)
	d, ok := data[0].Data.(invocationDataModel)
	if !ok {
		log.Println("Can not convert to invocationDataModel")
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
	assert.Nil(t, data.AgentVersion)
	assert.Nil(t, data.DataModelVersion)
	assert.Nil(t, data.ApplicationID)
	assert.Nil(t, data.ApplicationDomainName)
	assert.Nil(t, data.ApplicationClassName)
	assert.Nil(t, data.ApplicationName)
	assert.Nil(t, data.ApplicationVersion)
	assert.Nil(t, data.ApplicationStage)
	assert.Nil(t, data.ApplicationRuntime)
	assert.Nil(t, data.ApplicationRuntimeVersion)
	assert.Nil(t, data.ApplicationTags)

	assert.Equal(t, "AWS Lambda", data.ApplicationPlatform)
	assert.Equal(t, application.FunctionRegion, data.FunctionRegion)
	assert.NotNil(t, data.Tags)

	test.CleanEnvironment()
	ClearTags()
}

func TestPrepareDataStaticFieldsCompositeDataDisabled(t *testing.T) {
	config.ReportRestCompositeDataEnabled = false
	test.PrepareEnvironment()
	i := New()
	data := i.prepareData(context.TODO())
	assert.NotNil(t, data.ID)
	assert.Equal(t, "Invocation", data.Type)
	assert.Equal(t, constants.AgentVersion, *data.AgentVersion)
	assert.Equal(t, constants.DataModelVersion, *data.DataModelVersion)
	assert.Equal(t, test.AppID, *data.ApplicationID)
	assert.Equal(t, application.ApplicationDomainName, *data.ApplicationDomainName)
	assert.Equal(t, application.ApplicationClassName, *data.ApplicationClassName)
	assert.Equal(t, application.ApplicationName, *data.ApplicationName)
	assert.Equal(t, application.ApplicationVersion, *data.ApplicationVersion)
	assert.Equal(t, application.ApplicationStage, *data.ApplicationStage)
	assert.Equal(t, application.ApplicationRuntime, *data.ApplicationRuntime)
	assert.Equal(t, application.ApplicationRuntimeVersion, *data.ApplicationRuntimeVersion)
	assert.NotNil(t, *data.ApplicationTags)

	assert.Equal(t, "AWS Lambda", data.ApplicationPlatform)
	assert.Equal(t, application.FunctionRegion, data.FunctionRegion)
	assert.NotNil(t, data.Tags)

	test.CleanEnvironment()
	ClearTags()
}

func TestInvocationTags(t *testing.T) {
	agentTags := map[string]interface{}{
		"boolKey":   true,
		"intKey":    37,
		"floatKey":  3.14,
		"stringKey": "foobar",
		"dictKey": map[string]string{
			"key1": "val1",
			"key2": "val2",
		},
	}

	userTags := map[string]interface{}{
		"boolKey":   false,
		"intKey":    73,
		"floatKey":  6.21,
		"stringKey": "barfoo",
		"dictKey": map[string]string{
			"key1": "val3",
			"key2": "val4",
		},
	}

	for k, v := range agentTags {
		SetAgentTag(k, v)
	}

	for k, v := range userTags {
		SetTag(k, v)
	}

	assert.Equal(t, len(invocationTags), len(agentTags))
	assert.Equal(t, invocationTags["boolKey"], agentTags["boolKey"])
	assert.Equal(t, invocationTags["intKey"], agentTags["intKey"])
	assert.Equal(t, invocationTags["floatKey"], agentTags["floatKey"])
	assert.Equal(t, invocationTags["stringKey"], agentTags["stringKey"])
	assert.Equal(t, invocationTags["dictKey"], agentTags["dictKey"])

	assert.Equal(t, len(userInvocationTags), len(userTags))
	assert.Equal(t, userInvocationTags["boolKey"], userTags["boolKey"])
	assert.Equal(t, userInvocationTags["intKey"], userTags["intKey"])
	assert.Equal(t, userInvocationTags["floatKey"], userTags["floatKey"])
	assert.Equal(t, userInvocationTags["stringKey"], userTags["stringKey"])
	assert.Equal(t, userInvocationTags["dictKey"], userTags["dictKey"])

	ClearTags()
}
