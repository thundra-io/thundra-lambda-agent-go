package plugin

type Data interface{}

// Message defines the structure that given dataformat follows by Thundra. In here data could be a trace, metric or log data.
type Message struct {
	Data              Data   `json:"data"`
	Type              string `json:"type"`
	ApiKey            string `json:"apiKey"`
	DataFormatVersion string `json:"dataFormatVersion"`
}
