package agent

import (
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/thundra-io/thundra-lambda-agent-go/plugin"
)

func isThundraDisabled() bool {
	env := os.Getenv(thundraLambdaDisable)
	disabled, err := strconv.ParseBool(env)
	if err != nil {
		if env != "" {
			fmt.Println(err, " thundra_lambda_disable is not a bool value. Thundra is enabled by default.")
		}
		return false
	}
	return disabled
}

type timeoutError struct{}

func (e timeoutError) Error() string {
	return fmt.Sprintf("Lambda is timed out")
}

// determineTimeoutMargin fetches thundraLambdaTimeoutMargin if it exist, if not returns default timegap value
func determineTimeoutMargin() time.Duration {
	t := os.Getenv(thundraLambdaTimeoutMargin)
	// environment variable is not set
	if t == "" {
		return time.Duration(DefaultTimeoutMargin) * time.Millisecond
	}

	i, err := strconv.ParseInt(t, 10, 32)

	// environment variable is not set in the correct format
	if err != nil {
		fmt.Printf("%v: %s should be set with an integer\n", err, thundraLambdaTimeoutMargin)
		return time.Duration(DefaultTimeoutMargin) * time.Millisecond
	}

	return time.Duration(i) * time.Millisecond
}

// determineWarmup determines which warmup value to use. if warmup is set from environment variable, returns that value.
// Otherwise returns true if it's enabled by builder's enableWarmup method. Default value is false.
func determineWarmup() bool {
	w := os.Getenv(thundraLambdaWarmupWarmupAware)
	b, err := strconv.ParseBool(w)
	if err != nil {
		if w != "" {
			fmt.Println(err, " thundra_lambda_warmup_warmupAware should be set with a boolean.")
		}
		return false
	}
	return b
}

// determineApiKey determines which apiKey to use. if apiKey is set from environment variable, returns that value.
// Otherwise returns the value from builder's setApiKey method. Panic if it's not set by neither.
func determineAPIKey() {
	k := os.Getenv(thundraAPIKey)
	if k == "" {
		// TODO remove panics just log
		fmt.Println("Error no APIKey in env variables")
	}

	// Set it globally
	plugin.APIKey = k
}
