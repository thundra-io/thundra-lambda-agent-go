package samplers

import (
	"github.com/thundra-io/thundra-lambda-agent-go/tracer"
)

type errorAwareSampler struct {
}

func (e *errorAwareSampler) IsSampled(message interface{}) bool {
	if message != nil {
		switch data := message.(type) {
		case *tracer.RawSpan:
			if data != nil && data.Tags != nil {
				return data.Tags["error"] != nil
			}
		}
	}
	return false
}

func NewErrorAwareSampler() Sampler {
	return &errorAwareSampler{}
}
