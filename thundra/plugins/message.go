package plugins

type Data interface{}

// MonitoringDataWrapper defines the structure that given dataformat follows by Thundra. In here data could be a trace, metric or log data.
type MonitoringDataWrapper struct {
	DataModelVersion string `json:"dataModelVersion"`
	Type             string `json:"type"`
	Data             Data   `json:"data"`
	ApiKey           string `json:"apiKey"`
	Compressed       bool   `json:"compressed"`
}
