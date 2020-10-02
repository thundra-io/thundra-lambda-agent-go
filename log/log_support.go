package log

import "github.com/thundra-io/thundra-lambda-agent-go/v2/samplers"

var _sampler samplers.Sampler

func GetSampler() samplers.Sampler {
	return _sampler
}

func SetSampler(sampler samplers.Sampler) {
	_sampler = sampler
}
