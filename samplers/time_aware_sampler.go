package samplers

import (
	"sync/atomic"
	"time"
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

func NewTimeAwareSampler(freq int64) Sampler {
	return &timeAwareSampler{timeFreq: freq}
}
