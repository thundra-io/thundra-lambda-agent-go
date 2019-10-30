package test

import (
	"sync/atomic"

	"github.com/aws/aws-lambda-go/lambdacontext"
	"github.com/stretchr/testify/mock"
	"github.com/thundra-io/thundra-lambda-agent-go/application"
	"github.com/thundra-io/thundra-lambda-agent-go/plugin"
)

const (
	ApplicationName       = "TestFunctionName"
	FunctionName          = "TestFunctionName"
	ApplicationID         = "aws:lambda:TestRegion:guest:TestFunctionName"
	LogStreamName         = "2018/01/01/[$LATEST]1234567890"
	ApplicationInstanceID = "1234567890"
	FunctionVersion       = "$Version"
	ApplicationStage      = "TestStage"
	Region                = "TestRegion"
	MemoryLimit           = 512
	LogGroupName          = "TestLogGroupName"
)

//MockReporter is used in tests for mock reporter
type MockReporter struct {
	mock.Mock
	MessageQueue []plugin.MonitoringDataWrapper
	ReportedFlag *uint32
}

func (r *MockReporter) Collect(messages []plugin.MonitoringDataWrapper) {
	r.MessageQueue = append(r.MessageQueue, messages...)
	r.Called(messages)
}

func (r *MockReporter) Report() {
	r.Called()
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
func NewMockReporter() *MockReporter {
	r := &MockReporter{
		ReportedFlag: new(uint32),
	}
	r.On("Report").Return()
	r.On("ClearData").Return()
	r.On("Collect", mock.Anything).Return()
	return r
}

func PrepareEnvironment() {
	lambdacontext.LogStreamName = LogStreamName
	application.ApplicationName = ApplicationName
	application.FunctionName = FunctionName
	application.ApplicationInstanceID = ApplicationInstanceID
	application.ApplicationID = ApplicationID
	application.ApplicationVersion = FunctionVersion
	application.ApplicationStage = ApplicationStage
	application.FunctionRegion = Region
	application.MemoryLimit = MemoryLimit
	application.LogGroupName = LogGroupName
	application.LogStreamName = LogStreamName
}

func CleanEnvironment() {
	application.ApplicationName = ""
	application.FunctionName = ""
	application.ApplicationID = ""
	application.ApplicationInstanceID = ""
	application.ApplicationVersion = ""
	application.ApplicationStage = ""
	application.FunctionRegion = ""
	application.MemoryLimit = 0
}
