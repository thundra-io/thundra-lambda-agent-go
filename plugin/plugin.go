package plugin

import (
	"context"
	"encoding/json"

	"github.com/thundra-io/thundra-lambda-agent-go/config"

	"github.com/thundra-io/thundra-lambda-agent-go/constants"
)

var TraceID string
var TransactionID string

// Plugin interface provides necessary methods for the plugins to be used in thundra agent
type Plugin interface {
	BeforeExecution(ctx context.Context, request json.RawMessage) context.Context
	AfterExecution(ctx context.Context, request json.RawMessage, response interface{}, err interface{}) []MonitoringDataWrapper
	IsEnabled() bool
	Order() uint8
}

type Data interface{}

// MonitoringDataWrapper defines the structure that given dataformat follows by Thundra. In here data could be a trace, metric or log data.
type MonitoringDataWrapper struct {
	DataModelVersion string `json:"dataModelVersion"`
	Type             string `json:"type"`
	Data             Data   `json:"data"`
	APIKey           string `json:"apiKey"`
	Compressed       bool   `json:"compressed"`
}

func WrapMonitoringData(data interface{}, dataType string) MonitoringDataWrapper {
	return MonitoringDataWrapper{
		DataModelVersion: constants.DataModelVersion,
		Type:             dataType,
		Data:             data,
		APIKey:           config.APIKey,
	}
}
