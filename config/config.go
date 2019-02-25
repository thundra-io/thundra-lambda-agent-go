package config

import (
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/thundra-io/thundra-lambda-agent-go/constants"
)

var ThundraDisabled bool
var TimeoutMargin time.Duration
var WarmupEnabled bool
var DebugEnabled bool
var APIKey string
var TrustAllCertificates bool

func init() {
	ThundraDisabled = isThundraDisabled()
	DebugEnabled = isThundraDebugEnabled()
	TimeoutMargin = determineTimeoutMargin()
	WarmupEnabled = determineWarmup()
	APIKey = determineAPIKey()
	TrustAllCertificates = trustAllCertificates()
}

func isThundraDisabled() bool {
	env := os.Getenv(constants.ThundraLambdaDisable)
	disabled, err := strconv.ParseBool(env)
	if err != nil {
		if env != "" {
			fmt.Println(err, " thundra_lambda_disable is not a bool value. Thundra is enabled by default.")
		}
		return false
	}
	return disabled
}

func determineTimeoutMargin() time.Duration {
	t := os.Getenv(constants.ThundraLambdaTimeoutMargin)
	// environment variable is not set
	if t == "" {
		return time.Duration(constants.DefaultTimeoutMargin) * time.Millisecond
	}

	i, err := strconv.ParseInt(t, 10, 32)

	// environment variable is not set in the correct format
	if err != nil {
		fmt.Printf("%v: %s should be set with an integer\n", err, constants.ThundraLambdaTimeoutMargin)
		return time.Duration(constants.DefaultTimeoutMargin) * time.Millisecond
	}

	return time.Duration(i) * time.Millisecond
}

func determineWarmup() bool {
	w := os.Getenv(constants.ThundraLambdaWarmupWarmupAware)
	b, err := strconv.ParseBool(w)
	if err != nil {
		if w != "" {
			fmt.Println(err, " thundra_lambda_warmup_warmupAware should be set with a boolean.")
		}
		return false
	}
	return b
}

func determineAPIKey() string {
	apiKey := os.Getenv(constants.ThundraAPIKey)
	if apiKey == "" {
		fmt.Println("Error no APIKey in env variables")
	}
	return apiKey
}

func isThundraDebugEnabled() bool {
	b, err := strconv.ParseBool(os.Getenv(constants.ThundraLambdaDebugEnable))
	if err != nil {
		return false
	}
	return b
}

func trustAllCertificates() bool {
	b, err := strconv.ParseBool(os.Getenv(constants.ThundraTrustAllCertificates))
	if err != nil {
		return false
	}
	return b
}
