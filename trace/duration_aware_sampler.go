package trace

import (
	"github.com/thundra-io/thundra-lambda-agent-go/samplers"
	"github.com/thundra-io/thundra-lambda-agent-go/tracer"
)

type durationAwareSampler struct {
	duration   int64
	longerThan bool
}

func (d *durationAwareSampler) IsSampled(message interface{}) bool {
	if message != nil {
		switch data := message.(type) {
		case *tracer.RawSpan:
			if data != nil {
				duration := data.Duration()
				if d.longerThan {
					return duration > d.duration
				}
				return duration < d.duration
			}
		}
	}
	return false
}

func NewDurationAwareSampler(duration int64, longerThan bool) samplers.Sampler {
	return &durationAwareSampler{duration: duration, longerThan: longerThan}
}
