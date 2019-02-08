package metric

import (
	"fmt"

	uuid "github.com/satori/go.uuid"
	"github.com/shirou/gopsutil/mem"
	"github.com/thundra-io/thundra-lambda-agent-go/application"
)

func prepareMemoryMetricsData(mp *metricPlugin, base metricDataModel) metricDataModel {
	base.ID = uuid.NewV4().String()
	base.MetricName = memoryMetric

	memInfo, err := mem.VirtualMemory()
	if err != nil {
		fmt.Println(err)
	}

	procMemInfo, err := proc.MemoryInfo()
	if err != nil {
		fmt.Println(err)
	}

	base.Metrics = map[string]interface{}{
		appUsedMemory: procMemInfo.RSS,
		appMaxMemory:  application.MemoryLimit * 1024 * 1024,
		sysUsedMemory: memInfo.Used,
		sysMaxMemory:  memInfo.Total,
	}

	return base
}
