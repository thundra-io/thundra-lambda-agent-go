package samplers

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/thundra-io/thundra-lambda-agent-go/config"
	"github.com/thundra-io/thundra-lambda-agent-go/constants"
)

func TestDefaultCountFreq(t *testing.T) {

	cas := NewCountAwareSampler()

	assert.Equal(t, int64(constants.DefaultSamplingCountFreq), cas.(*countAwareSampler).countFreq)
}

func TestCountFreqFromEnv(t *testing.T) {
	config.SamplingCountFrequency = 10
	cas := NewCountAwareSampler()

	assert.Equal(t, int64(config.SamplingCountFrequency), cas.(*countAwareSampler).countFreq)
}

func TestFreqFromParam(t *testing.T) {
	cas := NewCountAwareSampler(10)

	assert.Equal(t, int64(10), cas.(*countAwareSampler).countFreq)
}

func TestSampledCountAware(t *testing.T) {
	cas := NewCountAwareSampler(2)

	assert.True(t, cas.IsSampled(nil))
	assert.False(t, cas.IsSampled(nil))
}
