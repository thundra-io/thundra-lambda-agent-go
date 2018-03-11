package thundra

import (
	"thundra-agent-go/plugin"
)

type Builder interface {
	AddPlugin(plugin.Plugin) Builder
	SetReporter(reporter) Builder
	Build() *thundra
}

type builder struct {
	plugins  []plugin.Plugin
	reporter reporter
}

func (b *builder) AddPlugin(plugin plugin.Plugin) Builder {
	b.plugins = append(b.plugins, plugin)
	return b
}

func (b *builder) SetReporter(reporter reporter) Builder {
	b.reporter = reporter
	return b
}

func (b *builder) Build() *thundra {
	if b.reporter == nil {
		b.reporter = &reporterImpl{}
	}
	return &thundra{
		b.plugins,
		b.reporter,
	}
}

func NewBuilder() Builder {
	return &builder{}
}
