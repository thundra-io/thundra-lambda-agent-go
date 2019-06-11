package trace

import "github.com/thundra-io/thundra-lambda-agent-go/samplers"

var _sampler samplers.Sampler

func GetSampler() samplers.Sampler {
	return _sampler
}

func SetSampler(sampler samplers.Sampler) {
	_sampler = sampler
}
