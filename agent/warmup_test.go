package agent

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCheckAndHandleWarmupNonEmptyPayload(t *testing.T) {
	payload := json.RawMessage(`{"firstName":"John","lastName":"Dow"}}`)
	assert.False(t, checkAndHandleWarmupRequest(payload))
}

func TestCheckAndHandleWarmupWarmCommand(t *testing.T) {
	payload := `#warmup wait=200`

	rawMessage, err := json.Marshal(payload)
	if err != nil {
		panic(err)
	}

	assert.True(t, checkAndHandleWarmupRequest(rawMessage))
}

func TestCheckAndHandleWarmupRequestEmptyPayload(t *testing.T) {
	payload := json.RawMessage(`{}`)
	assert.True(t, checkAndHandleWarmupRequest(payload))
}
