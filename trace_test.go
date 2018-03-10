package thundra

import (
	"testing"
	"context"
	"fmt"
	"encoding/json"
	"os"
	"time"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/aws/aws-lambda-go/lambdacontext"
)

const (
	duration           = 500
	functionName       = "TestFunctionName"
	memoryLimit        = 512
	functionVersion    = "$Version"
	applicationProfile = "TestProfile"
	apiKey             = "TestApiKey"
	defaultRegion      = "TestRegion"
)

func TestTrace(t *testing.T) {
	hello := func(s string) string {
		time.Sleep(time.Millisecond * duration)
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

			invocationStartTime := time.Now().Round(time.Millisecond)
			response, err := lambdaHandler.InvokeWithoutSerialization(context.TODO(), []byte(testCase.input))
			invocationEndTime := time.Now().Round(time.Millisecond)

			if testCase.expected.err != nil {
				assert.Equal(t, testCase.expected.err, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, testCase.expected.val, response.(string))

				//Monitor Data
				assert.Equal(t, traceDataType, mc.msg.Type)
				assert.Equal(t, apiKey, mc.msg.ApiKey)
				assert.Equal(t, dataFormatVersion, mc.msg.DataFormatVersion)

				//Trace Data
				assert.NotNil(t, mc.msg.Data.Id)
				assert.Equal(t, functionName, mc.msg.Data.ApplicationName)
				assert.Equal(t, "", mc.msg.Data.ApplicationId)
				assert.Equal(t, functionVersion, mc.msg.Data.ApplicationVersion)
				assert.Equal(t, applicationProfile, mc.msg.Data.ApplicationProfile)
				assert.Equal(t, applicationType, mc.msg.Data.ApplicationType)
				assert.NotNil(t, mc.msg.Data.ContextId)
				assert.Equal(t, functionName, mc.msg.Data.ContextName)
				assert.Equal(t, executionContext, mc.msg.Data.ContextType)

				st, err1 := time.Parse(timeFormat, mc.msg.Data.StartTime)
				et, err2 := time.Parse(timeFormat, mc.msg.Data.EndTime)
				if err1 != nil || err2 != nil {
					fmt.Println("err1: ", err1)
					fmt.Println("err2: ", err2)
				}
				assert.True(t, invocationStartTime.Before(st) || invocationStartTime.Equal(st))
				assert.True(t, st.Before(et))
				assert.True(t, et.Before(invocationEndTime))
				assert.True(t, int64(duration) <= mc.msg.Data.Duration)

				assert.Nil(t, mc.msg.Data.Errors)
				assert.Nil(t, mc.msg.Data.ThrownError)
				assert.Nil(t, mc.msg.Data.ThrownErrorMessage)

				//Trace Audit Info
				assert.Equal(t, functionName, mc.msg.Data.AuditInfo[auditInfoContextName])
				assert.NotNil(t, mc.msg.Data.AuditInfo[auditInfoId])
				assert.Equal(t, mc.msg.Data.StartTime, mc.msg.Data.AuditInfo[auditInfoOpenTime])
				assert.Equal(t, mc.msg.Data.EndTime, mc.msg.Data.AuditInfo[audit_info_close_time])
				assert.Nil(t, mc.msg.Data.AuditInfo[auditInfoErrors])
				assert.Nil(t, mc.msg.Data.AuditInfo[auditInfoThrownError])

				req := json.RawMessage(testCase.input)

				//Trace Properties
				assert.Equal(t, req, mc.msg.Data.Properties[auditInfoPropertiesRequest])
				assert.Equal(t, testCase.expected.val, mc.msg.Data.Properties[auditInfoPropertiesResponse])
				//TODO check cold start
				assert.Equal(t, "true", mc.msg.Data.Properties[auditInfoPropertiesColdStart])
				assert.Equal(t, defaultRegion, mc.msg.Data.Properties[auditInfoPropertiesFunctionRegion])
				assert.Equal(t, memoryLimit, mc.msg.Data.Properties[auditInfoPropertiesFunctionMemoryLimit])
			}
			cleanEnvironment()
		})
	}
}

func prepareEnvironment() {
	lambdacontext.FunctionName = functionName
	lambdacontext.MemoryLimitInMB = memoryLimit
	lambdacontext.FunctionVersion = functionVersion
	os.Setenv(ThundraApplicationProfile, applicationProfile)
	os.Setenv(ThundraApiKey, apiKey)
	os.Setenv(awsDefaultRegion, defaultRegion)
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
