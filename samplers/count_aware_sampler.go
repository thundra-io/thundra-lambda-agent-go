package samplers

import (
	"sync/atomic"

	"github.com/thundra-io/thundra-lambda-agent-go/config"
	"github.com/thundra-io/thundra-lambda-agent-go/constants"
)

type countAwareSampler struct {
	countFreq int64
	counter   int64
}

func (c *countAwareSampler) IsSampled(interface{}) bool {
	counter := atomic.AddInt64(&c.counter, 1)
	return (counter % c.countFreq) == 0
}

func NewCountAwareSampler(params ...int64) Sampler {
	var freq int64

	if config.SamplingCountFrequency > 0 {
		freq = int64(config.SamplingCountFrequency)
	} else if len(params) > 0 {
		freq = params[0]
	} else {
		freq = int64(constants.DefaultSamplingCountFreq)
	}

	return &countAwareSampler{countFreq: freq, counter: -1}
}
