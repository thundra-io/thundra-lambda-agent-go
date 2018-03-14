package plugin

type Data interface{}

//Data is TraceData
type Message struct {
	Data              Data   `json:"data"`
	Type              string `json:"type"`
	ApiKey            string `json:"apiKey"`
	DataFormatVersion string `json:"dataFormatVersion"`
}