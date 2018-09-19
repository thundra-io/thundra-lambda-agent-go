package metric

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/thundra-io/thundra-lambda-agent-go/plugin"
)

const numGoroutines = 5

//There are 2 goroutines running as default on testing
//One is the main and the other one is the testing
const defaultGoroutines = 2

func TestPrepareGoroutineStatsData(t *testing.T) {
	metric := NewBuilder().Build()
	metric.span.metricTimestamp = plugin.GetTimestamp()
	metric.span.startGCCount = 1
	metric.span.endGCCount = 2

	done := make(chan bool)
	generateGoroutines(done, numGoroutines)

	gcStatsData := prepareGoRoutineMetricsData(metric)

	assert.Equal(t, goroutineMetric, gcStatsData.MetricName)
	assert.Equal(t, metric.span.metricTimestamp, gcStatsData.MetricTimestamp)

	assert.Equal(t, uint64(numGoroutines+defaultGoroutines), gcStatsData.Metrics[numGoroutine])
	killGeneratedGoroutines(done, numGoroutines)
}

//Generates a number of Goroutines and wait for done signal
func generateGoroutines(done chan bool, numGoroutines int) {
	for i := 0; i < numGoroutines; i++ {
		go func(done chan bool) {
			<-done
		}(done)
	}
}

//Finished waiting goroutines
func killGeneratedGoroutines(done chan bool, numGoroutines int) {
	for i := 0; i < numGoroutines; i++ {
		done <- true
	}
}
