package samplers

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/thundra-io/thundra-lambda-agent-go/v2/config"
	"github.com/thundra-io/thundra-lambda-agent-go/v2/constants"
)

func TestDefaultTimeFreq(t *testing.T) {

	tas := NewTimeAwareSampler()

	assert.Equal(t, int64(constants.DefaultSamplingTimeFreq), tas.(*timeAwareSampler).timeFreq)
}

func TestTimeFreqFromEnv(t *testing.T) {
	config.SamplingTimeFrequency = 10
	tas := NewTimeAwareSampler()

	assert.Equal(t, int64(config.SamplingTimeFrequency), tas.(*timeAwareSampler).timeFreq)
}

func TestTimeFreqFromParam(t *testing.T) {
	config.SamplingTimeFrequency = -1
	tas := NewTimeAwareSampler(10)

	assert.Equal(t, int64(10), tas.(*timeAwareSampler).timeFreq)
}

func TestSampledTimeAware(t *testing.T) {
	tas := NewTimeAwareSampler(1)

	assert.True(t, tas.IsSampled(nil))
	time.Sleep(2000000)
	assert.True(t, tas.IsSampled(nil))
}
