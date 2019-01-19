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

func TestPrepareGoroutineMetricsData(t *testing.T) {
	mp := New()
	mp.data.metricTimestamp = plugin.GetTimestamp()
	mp.data.startGCCount = 1
	mp.data.endGCCount = 2

	done := make(chan bool)
	generateGoroutines(done, numGoroutines)
	base := mp.prepareMetricsData()
	grMetric := prepareGoRoutineMetricsData(mp, base)

	assert.True(t, len(grMetric.ID) != 0)
	assert.Equal(t, goroutineMetric, grMetric.MetricName)
	assert.Equal(t, mp.data.metricTimestamp, grMetric.MetricTimestamp)

	assert.Equal(t, uint64(numGoroutines+defaultGoroutines), grMetric.Metrics[numGoroutine])
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
