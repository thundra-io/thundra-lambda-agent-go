package metric

type metricData struct {
	//Base fields
	Id                        string                 `json:"id"`
	Type                      string                 `json:"type"`
	AgentVersion              string                 `json:"agentVersion"`
	DataModelVersion          string                 `json:"dataModelVersion"`
	ApplicationId             string                 `json:"applicationId"`
	ApplicationDomainName     string                 `json:"applicationDomainName"`
	ApplicationClassName      string                 `json:"applicationClassName"`
	ApplicationName           string                 `json:"applicationName"`
	ApplicationVersion        string                 `json:"applicationVersion"`
	ApplicationStage          string                 `json:"applicationStage"`
	ApplicationRuntime        string                 `json:"applicationRuntime"`
	ApplicationRuntimeVersion string                 `json:"applicationRuntimeVersion"`
	ApplicationTags           map[string]interface{} `json:"applicationTags"`

	TraceId         string                 `json:"traceId"`
	TracnsactionId  string                 `json:"transactionId"`
	SpanId          string                 `json:"spanId"`
	MetricName      string                 `json:"metricName"`
	MetricTimestamp int64                  `json:"metricTimestamp"`
	Metrics         map[string]interface{} `json:"metrics"`
	Tags            map[string]interface{} `json:"tags"`
}
