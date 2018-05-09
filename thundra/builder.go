package thundra

import (
	"os"

	"github.com/thundra-io/thundra-lambda-agent-go/plugin"
	"errors"
)

type Builder interface {
	AddPlugin(plugin.Plugin) Builder
	SetReporter(Reporter) Builder
	SetAPIKey(string) Builder
	EnableWarmup() Builder
	Build() *thundra
}

type builder struct {
	plugins  []plugin.Plugin
	reporter Reporter
	apiKey   string
	warmup   bool
}

func (b *builder) AddPlugin(plugin plugin.Plugin) Builder {
	b.plugins = append(b.plugins, plugin)
	return b
}

func (b *builder) SetReporter(reporter Reporter) Builder {
	b.reporter = reporter
	return b
}

func (b *builder) SetAPIKey(apiKey string) Builder {
	b.apiKey = apiKey
	return b
}

func (b *builder) EnableWarmup() Builder {
	b.warmup = true
	return b
}

func (b *builder) Build() *thundra {
	if b.reporter == nil {
		b.reporter = &reporterImpl{}
	}
	if b.apiKey == "" {
		k := os.Getenv(thundraApiKey)
		if k == "" {
			panic(errors.New("thundraApiKey is not set"))
		}
		b.apiKey = k
	}
	return &thundra{
		plugins:  b.plugins,
		reporter: b.reporter,
		apiKey:   b.apiKey,
		warmup:   b.warmup,
	}
}

func NewBuilder() Builder {
	return &builder{}
}
