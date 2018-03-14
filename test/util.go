package test

import "github.com/stretchr/testify/mock"

//MockReporter is used in tests for mock reporter
type MockReporter struct {
	mock.Mock
	MessageQueue []interface{}
}

func (r *MockReporter) Collect(msg interface{}) {
	r.MessageQueue = append(r.MessageQueue, msg)
	r.Called(msg)
}

func (r *MockReporter) Report(apiKey string) {
	r.Called(apiKey)
}

func (r *MockReporter) Clear() {
	r.Called()
}