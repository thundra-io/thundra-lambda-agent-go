package test

import (
	"sync/atomic"

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
	ReportedFlag *uint32
}

func (r *MockReporter) Collect(messages []interface{}) {
	r.MessageQueue = append(r.MessageQueue, messages...)
	r.Called(messages)
}

func (r *MockReporter) Report(apiKey string) {
	r.Called(apiKey)
	atomic.CompareAndSwapUint32(r.ReportedFlag, 0, 1)
}

func (r *MockReporter) ClearData() {
	r.Called()
}

func (r *MockReporter) Reported() *uint32 {
	return r.ReportedFlag
}

func (r *MockReporter) FlushFlag() {
	atomic.CompareAndSwapUint32(r.Reported(), 1, 0)
}

// NewMockReporter returns a new MockReporter
func NewMockReporter(testApiKey string) *MockReporter {
	r := &MockReporter{
		ReportedFlag: new(uint32),
	}
	r.On("Report", testApiKey).Return()
	r.On("ClearData").Return()
	r.On("Collect", mock.Anything).Return()
	return r
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
