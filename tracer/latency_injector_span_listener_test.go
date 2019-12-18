package tracer

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewLatencyInjectorFromConfig(t *testing.T) {
	config := map[string]interface{}{
		"delay":          float64(370),
		"injectOnFinish": true,
		"randomizeDelay": true,
		"addInfoTags":    false,
	}

	lsl := NewLatencyInjectorSpanListener(config).(*LatencyInjectorSpanListener)

	assert.Equal(t, int64(370), lsl.Delay)
	assert.Equal(t, true, lsl.InjectOnFinish)
	assert.Equal(t, true, lsl.RandomizeDelay)
	assert.Equal(t, false, lsl.AddInfoTags)

}

func TestNewLatencyInjectorFromConfigWithTypeErrors(t *testing.T) {
	config := map[string]interface{}{
		"injectOnFinish": 37,
		"delay":          "foo",
	}

	lsl := NewLatencyInjectorSpanListener(config).(*LatencyInjectorSpanListener)

	assert.Equal(t, false, lsl.InjectOnFinish)
	assert.Equal(t, defaultDelay, lsl.Delay)
	assert.Equal(t, false, lsl.RandomizeDelay)
	assert.Equal(t, true, lsl.AddInfoTags)
}
