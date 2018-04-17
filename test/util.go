package test

import "github.com/stretchr/testify/mock"

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