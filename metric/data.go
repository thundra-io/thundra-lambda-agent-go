package metric

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

	TraceID         string                 `json:"traceId"`
	TransactionID  string                 `json:"transactionId"`
	SpanID          string                 `json:"spanId"`
	MetricName      string                 `json:"metricName"`
	MetricTimestamp int64                  `json:"metricTimestamp"`
	Metrics         map[string]interface{} `json:"metrics"`
	Tags            map[string]interface{} `json:"tags"`
}
