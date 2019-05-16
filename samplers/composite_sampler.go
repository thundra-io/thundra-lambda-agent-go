package samplers

type compositeSampler struct {
	samplers []Sampler
	operator string
}

var defaultOperator = "or"

func (c *compositeSampler) IsSampled(data interface{}) bool {
	if c.samplers == nil || len(c.samplers) == 0 {
		return false
	}
	sampled := false
	if c.operator == "or" {
		for _, sampler := range c.samplers {
			sampled = sampler.IsSampled(data) || sampled
		}
		return sampled
	} else if c.operator == "and" {
		sampled = true
		for _, sampler := range c.samplers {
			sampled = sampler.IsSampled(data) && sampled
		}
		return sampled
	}

	return sampled
}

func NewCompositeSampler(samplers []Sampler, operator string) Sampler {
	_operator := operator
	if operator != "or" && operator != "and" {
		_operator = defaultOperator
	}
	return &compositeSampler{samplers, _operator}
}
