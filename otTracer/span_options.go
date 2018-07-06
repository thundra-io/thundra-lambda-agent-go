package otTracer

import (
	"github.com/opentracing/opentracing-go"
)

const (
	executionGroup   = "executionGroup"
	executionContext = "executionContext"
)

type startSpanOptions struct {
	Options opentracing.StartSpanOptions

	// operationGroup is lambda execution group
	operationGroup operationGroup

	// operationType is lambda execution type
	operationType operationType
}

type operationGroup string

// Apply satisfies the StartSpanOption interface.
func (r operationGroup) Apply(o *startSpanOptions) {
	o.operationGroup = executionGroup
}

type operationType string

// Apply satisfies the StartSpanOption interface.
func (r operationType) Apply(o *startSpanOptions) {
	o.operationType = executionContext
}

func newStartSpanOptions(sso []opentracing.StartSpanOption) startSpanOptions {
	opts := startSpanOptions{}
	for _, o := range sso {
		o.Apply(&opts.Options)
	}
	return opts
}