package metric

import (
	"github.com/thundra-io/thundra-lambda-agent-go/application"
	"github.com/thundra-io/thundra-lambda-agent-go/constants"
	"github.com/thundra-io/thundra-lambda-agent-go/plugin"
)

type metricDataModel struct {
	//Base fields
	ID                        string                 `json:"id"`
	Type                      string                 `json:"type"`
	AgentVersion              string                 `json:"agentVersion"`
	DataModelVersion          string                 `json:"dataModelVersion"`
	ApplicationID             string                 `json:"applicationId"`
	ApplicationDomainName     string                 `json:"applicationDomainName"`
	ApplicationClassName      string                 `json:"applicationClassName"`
	ApplicationName           string                 `json:"applicationName"`
	ApplicationVersion        string                 `json:"applicationVersion"`
	ApplicationStage          string                 `json:"applicationStage"`
	ApplicationRuntime        string                 `json:"applicationRuntime"`
	ApplicationRuntimeVersion string                 `json:"applicationRuntimeVersion"`
	ApplicationTags           map[string]interface{} `json:"applicationTags"`
	TraceID                   string                 `json:"traceId"`
	TransactionID             string                 `json:"transactionId"`
	SpanID                    string                 `json:"spanId"`
	MetricName                string                 `json:"metricName"`
	MetricTimestamp           int64                  `json:"metricTimestamp"`
	Metrics                   map[string]interface{} `json:"metrics"`
	Tags                      map[string]interface{} `json:"tags"`
}

func (mp *metricPlugin) prepareMetricsData() metricDataModel {
	return metricDataModel{
		Type:                      metricType,
		AgentVersion:              constants.AgentVersion,
		DataModelVersion:          constants.DataModelVersion,
		ApplicationID:             application.ApplicationID,
		ApplicationDomainName:     application.ApplicationDomainName,
		ApplicationClassName:      application.ApplicationClassName,
		ApplicationName:           application.ApplicationName,
		ApplicationVersion:        application.ApplicationVersion,
		ApplicationStage:          application.ApplicationStage,
		ApplicationRuntime:        application.ApplicationRuntime,
		ApplicationRuntimeVersion: application.ApplicationRuntimeVersion,
		ApplicationTags:           application.ApplicationTags,
		TraceID:                   plugin.TraceID,
		TransactionID:             plugin.TransactionID,
		SpanID:                    "", // Optional
		MetricTimestamp:           mp.data.metricTimestamp,
		Metrics:                   map[string]interface{}{},
		Tags:                      map[string]interface{}{},
	}
}
