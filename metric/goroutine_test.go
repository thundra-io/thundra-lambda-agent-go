package metric

import (
	"testing"
	"github.com/thundra-io/thundra-lambda-agent-go/plugin"
	"github.com/stretchr/testify/assert"
)

const numGoroutines = 5

//There are 2 goroutines running as default on testing
//One is the main and the other one is the testing
const defaultGoroutines = 2

func TestPrepareGoroutineStatsData(t *testing.T) {
	metric := NewBuilder().EnableGoroutineStats().Build()
	metric.statTimestamp = plugin.GetTimestamp()
	metric.startGCCount = 1
	metric.endGCCount = 2

	done := make(chan bool)
	generateGoroutines(done, numGoroutines)

	gcStatsData := prepareGoRoutineStatsData(metric)

	assert.Equal(t, goroutineStat, gcStatsData.StatName)
	assert.Equal(t, metric.statTimestamp, gcStatsData.StatTimestamp)

	assert.Equal(t, uint64(numGoroutines+defaultGoroutines), gcStatsData.NumGoroutine)
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
