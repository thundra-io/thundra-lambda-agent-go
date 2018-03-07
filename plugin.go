package thundra

import (
	"context"
	"sync"
	"encoding/json"
)

type Plugin interface {
	BeforeExecution(ctx context.Context, request interface{}, wg *sync.WaitGroup)
	AfterExecution(ctx context.Context, request interface{}, response interface{}, error interface{}, wg *sync.WaitGroup)
	OnPanic(ctx context.Context, request json.RawMessage, panic *ThundraPanic, wg *sync.WaitGroup)
}

type CollecterAwarePlugin interface {
	SetCollector(collector *collector)
}

type PluginFactory interface {
	Create() Plugin
}

var pluginDictionary map[string]PluginFactory

func discoverPlugins() {
	pD := make(map[string]PluginFactory)
	//TODO read plugin list from file
	pD["trace"] = &TraceFactory{}
	pluginDictionary = pD
}

func registerPluginFactory(pluginName string, factory PluginFactory){
	pluginDictionary[pluginName] = factory
}