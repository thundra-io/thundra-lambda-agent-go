package samplers

import (
	"sync/atomic"
	"time"

	"github.com/thundra-io/thundra-lambda-agent-go/config"
	"github.com/thundra-io/thundra-lambda-agent-go/constants"
)

type timeAwareSampler struct {
	timeFreq   int64
	latestTime int64
}

func (t *timeAwareSampler) IsSampled(interface{}) bool {
	sampled := false
	now := time.Now().UnixNano() / (int64(time.Millisecond) / int64(time.Nanosecond))
	latestTime := atomic.LoadInt64(&t.latestTime)
	atomic.StoreInt64(&t.latestTime, now)
	if now > latestTime+t.timeFreq {
		sampled = true
	}
	return sampled
}

func NewTimeAwareSampler(params ...int64) Sampler {
	var freq int64

	if config.SamplingTimeFrequency > 0 {
		freq = int64(config.SamplingTimeFrequency)
	} else if len(params) > 0 {
		freq = params[0]
	} else {
		freq = int64(constants.DefaultSamplingTimeFreq)
	}

	return &timeAwareSampler{timeFreq: freq}
}
