package plugin

import (
	"context"
	"encoding/json"
)

// Plugin interface provides necessary methods for the plugins to be used in thundra agent
type Plugin interface {
	BeforeExecution(ctx context.Context, request json.RawMessage) context.Context
	AfterExecution(ctx context.Context, request json.RawMessage, response interface{}, err interface{}) []MonitoringDataWrapper
	IsEnabled() bool
}
