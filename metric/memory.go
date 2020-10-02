package metric

import (
	"log"

	uuid "github.com/google/uuid"
	"github.com/thundra-io/thundra-lambda-agent-go/v2/application"
)

const miBToB = 1024 * 1024

func prepareMemoryMetricsData(mp *metricPlugin, base metricDataModel) metricDataModel {
	base.ID = uuid.New().String()
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
