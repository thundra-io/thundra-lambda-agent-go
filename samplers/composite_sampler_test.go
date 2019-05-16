package samplers

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

type mockSampler struct {
	sampled bool
}

func (c *mockSampler) IsSampled(data interface{}) bool {
	return c.sampled
}

func newMockSampler(sampled bool) Sampler {
	return &mockSampler{sampled}
}

func TestWithNoSamplers(t *testing.T) {
	cms := NewCompositeSampler(nil, "")
	assert.False(t, cms.IsSampled(nil))

	var samplers []Sampler
	cms = NewCompositeSampler(samplers, "")
}

func TestAndOperator(t *testing.T) {
	samplers := []Sampler{newMockSampler(true), newMockSampler(true)}

	cms := NewCompositeSampler(samplers, "and")
	assert.True(t, cms.IsSampled(nil))

	samplers = append(samplers, newMockSampler(false))

	cms = NewCompositeSampler(samplers, "and")
	assert.False(t, cms.IsSampled(nil))
}

func TestOrOperator(t *testing.T) {
	samplers := []Sampler{newMockSampler(false), newMockSampler(false)}

	cms := NewCompositeSampler(samplers, "or")
	assert.False(t, cms.IsSampled(nil))

	samplers = append(samplers, newMockSampler(true))

	cms = NewCompositeSampler(samplers, "or")
	assert.True(t, cms.IsSampled(nil))
}

func TestDefaultOperator(t *testing.T) {
	samplers := []Sampler{newMockSampler(false), newMockSampler(false)}

	cms := NewCompositeSampler(samplers, "")
	assert.False(t, cms.IsSampled(nil))

	samplers = append(samplers, newMockSampler(true))

	cms = NewCompositeSampler(samplers, "")
	assert.True(t, cms.IsSampled(nil))
}
