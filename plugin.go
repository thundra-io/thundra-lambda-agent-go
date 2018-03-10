package thundra

import (
	"context"
	"sync"
	"encoding/json"
)

type Plugin interface {
	BeforeExecution(ctx context.Context, request interface{}, wg *sync.WaitGroup)
	AfterExecution(ctx context.Context, request interface{}, response interface{}, error interface{}, wg *sync.WaitGroup) Message
	OnPanic(ctx context.Context, request json.RawMessage, panic *ThundraPanic, wg *sync.WaitGroup) Message
}

type CollecterAwarePlugin interface {
	Plugin
	SetCollector(collector collector)
}

type PluginFactory interface {
	Create() Plugin
}