package test

import (
	"github.com/aws/aws-lambda-go/lambdacontext"
	"github.com/stretchr/testify/mock"
	"github.com/thundra-io/thundra-lambda-agent-go/plugin"
)

const (
	FunctionName       = "TestFunctionName"
	LogStreamName      = "2018/01/01/[$LATEST]1234567890"
	AppId              = "1234567890"
	FunctionVersion    = "$Version"
	ApplicationProfile = "TestProfile"
	Region             = "TestRegion"
	MemoryLimit        = 512
	LogGroupName       = "TestLogGroupName"
)

//MockReporter is used in tests for mock reporter
type MockReporter struct {
	mock.Mock
	MessageQueue []interface{}
}

func (r *MockReporter) Collect(messages []interface{}) {
	r.MessageQueue = append(r.MessageQueue, messages...)
	r.Called(messages)
}

func (r *MockReporter) Report(apiKey string) {
	r.Called(apiKey)
}

func (r *MockReporter) Clear() {
	r.Called()
}

func PrepareEnvironment() {
	lambdacontext.LogStreamName = LogStreamName
	plugin.ApplicationName = FunctionName
	plugin.ApplicationId = AppId
	plugin.ApplicationVersion = FunctionVersion
	plugin.ApplicationProfile = ApplicationProfile
	plugin.Region = Region
	plugin.MemorySize = MemoryLimit
	plugin.LogGroupName = LogGroupName
	plugin.LogStreamName = LogStreamName
}

func CleanEnvironment() {
	plugin.ApplicationName = ""
	plugin.ApplicationId = ""
	plugin.ApplicationVersion = ""
	plugin.ApplicationProfile = ""
	plugin.Region = ""
	plugin.MemorySize = 0
}
