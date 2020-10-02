package metric

import (
	"github.com/thundra-io/thundra-lambda-agent-go/v2/application"
	"github.com/thundra-io/thundra-lambda-agent-go/v2/plugin"
)

type metricDataModel struct {
	//Base fields
	plugin.BaseDataModel
	ID              string                 `json:"id"`
	Type            string                 `json:"type"`
	TraceID         string                 `json:"traceId"`
	TransactionID   string                 `json:"transactionId"`
	SpanID          string                 `json:"spanId"`
	MetricName      string                 `json:"metricName"`
	MetricTimestamp int64                  `json:"metricTimestamp"`
	Metrics         map[string]interface{} `json:"metrics"`
	Tags            map[string]interface{} `json:"tags"`
}

func (mp *metricPlugin) prepareMetricsData() metricDataModel {
	return metricDataModel{
		BaseDataModel:   plugin.GetBaseData(),
		Type:            metricType,
		TraceID:         plugin.TraceID,
		TransactionID:   plugin.TransactionID,
		SpanID:          "", // Optional
		MetricTimestamp: mp.data.metricTimestamp,
		Metrics:         map[string]interface{}{},
		Tags: map[string]interface{}{
			"aws.region": application.FunctionRegion,
		},
	}
}
