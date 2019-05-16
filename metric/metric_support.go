package metric

import "github.com/thundra-io/thundra-lambda-agent-go/samplers"

var _sampler = samplers.NewCompositeSampler([]samplers.Sampler{samplers.NewTimeAwareSampler(), samplers.NewCountAwareSampler()}, "or")

func GetSampler() samplers.Sampler {
	return _sampler
}

func SetSampler(sampler samplers.Sampler) {
	_sampler = sampler
}
