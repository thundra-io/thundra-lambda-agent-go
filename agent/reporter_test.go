package agent

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/thundra-io/thundra-lambda-agent-go/config"
	"github.com/thundra-io/thundra-lambda-agent-go/plugin"
	"github.com/thundra-io/thundra-lambda-agent-go/test"
)

type RoundTripFunc func(req *http.Request) (*http.Response, error)

func (f RoundTripFunc) RoundTrip(req *http.Request) (*http.Response, error) {
	return f(req)
}

func newTestClient(fn RoundTripFunc) *http.Client {
	return &http.Client{
		Transport: RoundTripFunc(fn),
	}
}

func newTestReporter(fn RoundTripFunc) *reporterImpl {
	return &reporterImpl{
		client:   newTestClient(fn),
		reported: new(uint32),
	}
}

type mockDataModel struct{}

var mockData mockDataModel

func TestCollect(t *testing.T) {
	test.PrepareEnvironment()
	messages := []plugin.MonitoringDataWrapper{plugin.WrapMonitoringData(mockData, "Invocation")}
	testReporter := newTestReporter(func(req *http.Request) (*http.Response, error) {
		return &(http.Response{}), nil
	})
	testReporter.Collect(messages)
	assert.Equal(t, messages, testReporter.messageQueue)
	test.CleanEnvironment()
}

func TestCollectAsyncCompositeDisabled(t *testing.T) {
	config.ReportPublishCloudwatchEnabled = true
	config.ReportCloudwatchCompositeDataEnabled = false
	test.PrepareEnvironment()
	messages := []plugin.MonitoringDataWrapper{plugin.WrapMonitoringData(mockData, "Invocation")}
	testReporter := newTestReporter(func(req *http.Request) (*http.Response, error) {
		return &(http.Response{}), nil
	})
	testReporter.Collect(messages)
	var expectedMessages []plugin.MonitoringDataWrapper
	assert.Equal(t, expectedMessages, testReporter.messageQueue)
	test.CleanEnvironment()
}

func TestClearData(t *testing.T) {
	test.PrepareEnvironment()
	messages := []plugin.MonitoringDataWrapper{plugin.WrapMonitoringData(mockData, "Invocation")}
	testReporter := newTestReporter(func(req *http.Request) (*http.Response, error) {
		return &(http.Response{}), nil
	})
	testReporter.messageQueue = messages
	testReporter.ClearData()
	assert.Equal(t, []plugin.MonitoringDataWrapper{}, testReporter.messageQueue)
	test.CleanEnvironment()
}

func TestReportComposite(t *testing.T) {
	test.PrepareEnvironment()
	messages := []plugin.MonitoringDataWrapper{plugin.WrapMonitoringData(mockData, "Invocation")}
	testReporter := newTestReporter(func(req *http.Request) (*http.Response, error) {
		body, _ := ioutil.ReadAll(req.Body)
		var data plugin.MonitoringDataWrapper
		json.Unmarshal(body, &data)
		assert.Equal(t, "Composite", data.Type)
		return &(http.Response{}), nil
	})
	testReporter.messageQueue = messages
	testReporter.Report()
	test.CleanEnvironment()
}

func TestReportCompositeDisabled(t *testing.T) {
	config.ReportRestCompositeDataEnabled = false
	test.PrepareEnvironment()
	messages := []plugin.MonitoringDataWrapper{plugin.WrapMonitoringData(mockData, "Invocation")}
	testReporter := newTestReporter(func(req *http.Request) (*http.Response, error) {
		body, _ := ioutil.ReadAll(req.Body)
		var data []plugin.MonitoringDataWrapper
		json.Unmarshal(body, &data)
		assert.Equal(t, "Invocation", data[0].Type)
		return &(http.Response{}), nil
	})
	testReporter.messageQueue = messages
	testReporter.Report()
	test.CleanEnvironment()
}
