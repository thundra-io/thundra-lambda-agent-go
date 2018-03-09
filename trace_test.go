package thundra

import (
	"testing"
	"context"
	"fmt"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"encoding/json"
	"github.com/aws/aws-lambda-go/lambdacontext"
	"os"
	"thundra-agent-go/constants"
	"time"
)

const DURATION = 500
const FUNCTION_NAME = "TestFunctionName"
const MEMORY_LIMIT = 512
const FUNCTION_VERSION = "$Version"
const APPLICATION_PROFILE = "TestProfile"
const API_KEY = "TestApiKey"
const DEFAULT_REGION = "TestRegion"

func TestTrace(t *testing.T) {
	hello := func(s string) string {
		time.Sleep(time.Millisecond * DURATION)
		return fmt.Sprintf("Happy monitoring with %s!", s)
	}

	testCases := []struct {
		name     string
		input    string
		expected expected
		handler  interface{}
	}{
		{
			input:    `"Thundra"`,
			expected: expected{"Happy monitoring with Thundra!", nil},
			handler: func(name string) (string, error) {
				return hello(name), nil
			},
		},
	}
	for i, testCase := range testCases {
		t.Run(fmt.Sprintf("testCase[%d] %s", i, testCase.name), func(t *testing.T) {
			prepareEnvironment()

			mc := new(MockCollector)
			mc.On("report").Return()
			mc.On("clear").Return()
			mc.On("collect", mock.Anything).Return()

			th := createNewWithCollector([]string{"trace"}, mc)
			lambdaHandler := WrapLambdaHandler(testCase.handler, th)

			invocationStartTime := time.Now()
			response, err := lambdaHandler.InvokeWithoutSerialization(context.TODO(), []byte(testCase.input))
			invocationEndTime := time.Now()

			if testCase.expected.err != nil {
				assert.Equal(t, testCase.expected.err, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, testCase.expected.val, response.(string))

				//Monitor Data
				assert.Equal(t, constants.DATA_TYPE, mc.msg.Type)
				assert.Equal(t, API_KEY, mc.msg.ApiKey)
				assert.Equal(t, constants.DATA_FORMAT_VERSION, mc.msg.DataFormatVersion)

				//Trace Data
				assert.NotNil(t, mc.msg.Data.Id)
				assert.Equal(t, FUNCTION_NAME, mc.msg.Data.ApplicationName)
				assert.Equal(t, "", mc.msg.Data.ApplicationId)
				assert.Equal(t, FUNCTION_VERSION, mc.msg.Data.ApplicationVersion)
				assert.Equal(t, APPLICATION_PROFILE, mc.msg.Data.ApplicationProfile)
				assert.Equal(t, constants.APPLICATION_TYPE, mc.msg.Data.ApplicationType)
				assert.NotNil(t, mc.msg.Data.ContextId)
				assert.Equal(t, FUNCTION_NAME, mc.msg.Data.ContextName)
				assert.Equal(t, constants.EXECUTION_CONTEXT, mc.msg.Data.ContextType)

				st, err1 := time.Parse(constants.TIME_FORMAT, mc.msg.Data.StartTime)
				et, err2 := time.Parse(constants.TIME_FORMAT, mc.msg.Data.EndTime)
				if err1 != nil || err2 != nil {
					fmt.Println("err1: ", err1)
					fmt.Println("err2: ", err2)
				}
				assert.True(t, invocationStartTime.Before(st) || invocationStartTime.Equal(st))
				assert.True(t, st.Before(et))
				assert.True(t, et.Before(invocationEndTime))
				assert.True(t, int64(DURATION) <= mc.msg.Data.Duration)

				assert.Nil(t, mc.msg.Data.Errors)
				assert.Nil(t, mc.msg.Data.ThrownError)
				assert.Nil(t, mc.msg.Data.ThrownErrorMessage)

				//Trace Audit Info
				assert.Equal(t, FUNCTION_NAME, mc.msg.Data.AuditInfo[constants.AUDIT_INFO_CONTEXT_NAME])
				assert.NotNil(t, mc.msg.Data.AuditInfo[constants.AUDIT_INFO_ID])
				assert.Equal(t, mc.msg.Data.StartTime, mc.msg.Data.AuditInfo[constants.AUDIT_INFO_OPEN_TIME])
				assert.Equal(t, mc.msg.Data.EndTime, mc.msg.Data.AuditInfo[constants.AUDIT_INFO_CLOSE_TIME])
				assert.Nil(t, mc.msg.Data.AuditInfo[constants.AUDIT_INFO_ERRORS])
				assert.Nil(t, mc.msg.Data.AuditInfo[constants.AUDIT_INFO_THROWN_ERROR])

				req := json.RawMessage(testCase.input)

				//Trace Properties
				assert.Equal(t, req, mc.msg.Data.Properties[constants.AUDIT_INFO_PROPERTIES_REQUEST])
				assert.Equal(t, testCase.expected.val, mc.msg.Data.Properties[constants.AUDIT_INFO_PROPERTIES_RESPONSE])
				//TODO check cold start
				assert.Equal(t, "true", mc.msg.Data.Properties[constants.AUDIT_INFO_PROPERTIES_COLD_START])
				assert.Equal(t, DEFAULT_REGION, mc.msg.Data.Properties[constants.AUDIT_INFO_PROPERTIES_FUNCTION_REGION])
				assert.Equal(t, MEMORY_LIMIT, mc.msg.Data.Properties[constants.AUDIT_INFO_PROPERTIES_FUNCTION_MEMORY_LIMIT])
			}
			cleanEnvironment()
		})
	}

}

func prepareEnvironment() {
	lambdacontext.FunctionName = FUNCTION_NAME
	lambdacontext.MemoryLimitInMB = MEMORY_LIMIT
	lambdacontext.FunctionVersion = FUNCTION_VERSION
	os.Setenv(constants.THUNDRA_APPLICATION_PROFILE, APPLICATION_PROFILE)
	os.Setenv(constants.THUNDRA_API_KEY, API_KEY)
	os.Setenv(constants.AWS_DEFAULT_REGION, DEFAULT_REGION)
}

func cleanEnvironment() {
	lambdacontext.FunctionName = ""
	lambdacontext.MemoryLimitInMB = 0
	lambdacontext.FunctionVersion = ""
	os.Clearenv()
}

func (handler thundraLambdaHandler) InvokeWithoutSerialization(ctx context.Context, payload []byte) (interface{}, error) {
	response, err := handler(ctx, payload)
	if err != nil {
		return nil, err
	}
	return response, nil
}
