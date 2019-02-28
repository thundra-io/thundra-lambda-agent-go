package config

import (
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/thundra-io/thundra-lambda-agent-go/constants"
)

var ThundraDisabled bool
var TraceDisabled bool
var MetricDisabled bool
var LogDisabled bool
var LogLevel string
var TraceRequestDisabled bool
var TraceResponseDisabled bool
var TimeoutMargin time.Duration
var WarmupEnabled bool
var DebugEnabled bool
var APIKey string
var TrustAllCertificates bool

func init() {
	ThundraDisabled = isThundraDisabled()
	TraceDisabled = isTraceDisabled()
	MetricDisabled = isMetricDisabled()
	LogDisabled = isLogDisabled()
	TraceRequestDisabled = isTraceRequestDisabled()
	TraceResponseDisabled = isTraceResponseDisabled()
	DebugEnabled = isThundraDebugEnabled()
	TimeoutMargin = determineTimeoutMargin()
	WarmupEnabled = determineWarmup()
	APIKey = determineAPIKey()
	LogLevel = determineLogLevel()
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

func isTraceDisabled() bool {
	env := os.Getenv(constants.ThundraDisableTrace)
	disabled, err := strconv.ParseBool(env)
	if err != nil {
		if env != "" {
			fmt.Println(err, constants.ThundraDisableTrace+" is not a bool value. Trace plugin is enabled by default.")
		}
		return false
	}
	return disabled
}

func isMetricDisabled() bool {
	env := os.Getenv(constants.ThundraDisableMetric)
	disabled, err := strconv.ParseBool(env)
	if err != nil {
		if env != "" {
			fmt.Println(err, constants.ThundraDisableMetric+" is not a bool value. Metric plugin is enabled by default.")
		}
		return false
	}
	return disabled
}

func isLogDisabled() bool {
	env := os.Getenv(constants.ThundraDisableLog)
	disabled, err := strconv.ParseBool(env)
	if err != nil {
		if env != "" {
			fmt.Println(err, constants.ThundraDisableLog+" is not a bool value. Log plugin is enabled by default.")
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

func isTraceRequestDisabled() bool {
	env := os.Getenv(constants.ThundraDisableTraceRequest)
	disabled, err := strconv.ParseBool(env)
	if err != nil {
		if env != "" {
			fmt.Println(err, constants.ThundraDisableTraceRequest+"is not a bool value. Trace request is not disabled.")
		}
		return false
	}
	return disabled
}

func isTraceResponseDisabled() bool {
	env := os.Getenv(constants.ThundraDisableTraceResponse)
	disabled, err := strconv.ParseBool(env)
	if err != nil {
		if env != "" {
			fmt.Println(err, constants.ThundraDisableTraceResponse+" is not a bool value. Trace response is not disabled.")
		}
		return false
	}
	return disabled
}

func determineLogLevel() string {
	level := os.Getenv(constants.ThundraLogLogLevel)
	return strings.ToUpper(level)
}
