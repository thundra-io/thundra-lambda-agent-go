package samplers

import (
	"github.com/thundra-io/thundra-lambda-agent-go/v2/tracer"
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

func NewDurationAwareSampler(duration int64, longerThanArr ...bool) Sampler {
	var longerThan bool
	if len(longerThanArr) > 0 {
		longerThan = longerThanArr[0]
	}
	return &durationAwareSampler{duration: duration, longerThan: longerThan}
}
