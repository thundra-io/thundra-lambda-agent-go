package metric

import (
	"log"

	uuid "github.com/satori/go.uuid"
	"github.com/thundra-io/thundra-lambda-agent-go/application"
)

const miBToB = 1024 * 1024

func prepareMemoryMetricsData(mp *metricPlugin, base metricDataModel) metricDataModel {
	base.ID = uuid.NewV4().String()
	base.MetricName = memoryMetric

	procMemInfo, err := proc.MemoryInfo()
	if err != nil {
		log.Println(err)
	}
	
	application.MemoryUsed = int(procMemInfo.RSS / miBToB)

	base.Metrics = map[string]interface{}{
		appUsedMemory: procMemInfo.RSS,
		appMaxMemory:  application.MemoryLimit * miBToB,
	}

	return base
}
