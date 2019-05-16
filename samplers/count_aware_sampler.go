package samplers

import (
	"sync/atomic"
)

type countAwareSampler struct {
	countFreq int64
	counter   int64
}

func (c *countAwareSampler) IsSampled(interface{}) bool {
	counter := atomic.AddInt64(&c.counter, 1)
	return (counter % c.countFreq) == 0
}

func NewCountAwareSampler(freq int64) Sampler {
	return &countAwareSampler{countFreq: freq, counter: -1}
}
