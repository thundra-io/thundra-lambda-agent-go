package config

import (
	"log"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/thundra-io/thundra-lambda-agent-go/constants"
)

var ThundraDisabled bool
var TraceDisabled bool
var MetricDisabled bool
var AwsIntegrationDisabled bool
var LogDisabled bool
var LogLevel string
var TraceRequestDisabled bool
var TraceResponseDisabled bool
var TimeoutMargin time.Duration
var WarmupEnabled bool
var DebugEnabled bool
var APIKey string
var TrustAllCertificates bool
var MaskDynamoDBStatement bool
var DynamoDBTraceInjectionEnabled bool
var LambdaTraceInjectionDisabled bool
var MaskRDBStatement bool
var MaskEsBody bool
var MaskRedisCommand bool
var MaskMongoDBCommand bool
var MaskSNSMessage bool
var MaskSQSMessage bool
var MaskLambdaPayload bool
var MaskHTTPBody bool
var MaskAthenaStatement bool

var TraceKinesisRequestEnabled bool
var TraceFirehoseRequestEnabled bool
var TraceCloudwatchlogRequestEnabled bool

var ReportRestCompositeBatchSize int
var ReportCloudwatchCompositeBatchSize int

var ReportRestCompositeDataEnabled bool
var ReportCloudwatchCompositeDataEnabled bool
var ReportCloudwatchEnabled bool

var SamplingCountFrequency int
var SamplingTimeFrequency int

var HTTPIntegrationUrlPathDepth int

var AwsLambdaFunctionMemorySize int
var AwsLambdaRegion string

func init() {
	ThundraDisabled = boolFromEnv(constants.ThundraLambdaDisable, false)
	TraceDisabled = boolFromEnv(constants.ThundraDisableTrace, false)
	MetricDisabled = boolFromEnv(constants.ThundraDisableMetric, false)
	LogDisabled = boolFromEnv(constants.ThundraDisableLog, false)
	AwsIntegrationDisabled = boolFromEnv(constants.ThundraDisableAwsIntegration, false)
	TraceRequestDisabled = boolFromEnv(constants.ThundraDisableTraceRequest, false)
	TraceResponseDisabled = boolFromEnv(constants.ThundraDisableTraceResponse, false)
	DebugEnabled = boolFromEnv(constants.ThundraLambdaDebugEnable, false)
	WarmupEnabled = boolFromEnv(constants.ThundraLambdaWarmupWarmupAware, false)
	APIKey = determineAPIKey()
	LogLevel = determineLogLevel()
	TrustAllCertificates = boolFromEnv(constants.ThundraTrustAllCertificates, false)
	MaskDynamoDBStatement = boolFromEnv(constants.ThundraMaskDynamoDBStatement, false)
	MaskAthenaStatement = boolFromEnv(constants.ThundraMaskAthenaStatement, false)
	MaskRDBStatement = boolFromEnv(constants.ThundraMaskRDBStatement, false)
	MaskEsBody = boolFromEnv(constants.ThundraMaskEsBody, false)
	MaskRedisCommand = boolFromEnv(constants.ThundraMaskRedisCommand, false)
	TraceKinesisRequestEnabled = boolFromEnv(constants.ThundraLambdaTraceKinesisRequestEnable, false)
	TraceFirehoseRequestEnabled = boolFromEnv(constants.ThundraLambdaTraceFirehoseRequestEnable, false)
	TraceCloudwatchlogRequestEnabled = boolFromEnv(constants.ThundraLambdaTraceCloudwatchlogRequestEnable, false)
	DynamoDBTraceInjectionEnabled = boolFromEnv(constants.EnableDynamoDbTraceInjection, false)
	LambdaTraceInjectionDisabled = boolFromEnv(constants.DisableLambdaTraceInjection, false)
	ReportCloudwatchCompositeBatchSize = intFromEnv(constants.ThundraLambdaReportCloudwatchCompositeBatchSize,
		constants.ThundraLambdaReportCloudwatchCompositeBatchSizeDefault)
	ReportRestCompositeBatchSize = intFromEnv(constants.ThundraLambdaReportRestCompositeBatchSize,
		constants.ThundraLambdaReportRestCompositeBatchSizeDefault)
	ReportRestCompositeDataEnabled = boolFromEnv(constants.ThundraLambdaReportRestCompositeEnable, true)
	ReportCloudwatchCompositeDataEnabled = boolFromEnv(constants.ThundraLambdaReportCloudwatchCompositeEnable, true)
	ReportCloudwatchEnabled = boolFromEnv(constants.ThundraLambdaReportCloudwatchEnable, false)
	MaskMongoDBCommand = boolFromEnv(constants.ThundraMaskMongoDBCommand, false)
	SamplingCountFrequency = intFromEnv(constants.ThundraAgentMetricCountAwareSamplerCountFreq, -1)
	SamplingTimeFrequency = intFromEnv(constants.ThundraAgentMetricTimeAwareSamplerTimeFreq, -1)
	MaskSNSMessage = boolFromEnv(constants.ThundraMaskSNSMessage, false)
	MaskSQSMessage = boolFromEnv(constants.ThundraMaskSQSMessage, false)
	MaskLambdaPayload = boolFromEnv(constants.ThundraMaskLambdaPayload, false)
	MaskHTTPBody = boolFromEnv(constants.ThundraMaskHTTPBody, false)
	HTTPIntegrationUrlPathDepth = intFromEnv(constants.ThundraAgentTraceIntegrationsHttpUrlDepth, 1)
	AwsLambdaFunctionMemorySize = intFromEnv(constants.AwsLambdaFunctionMemorySize, -1)
	AwsLambdaRegion = os.Getenv(constants.AwsLambdaRegion)
	TimeoutMargin = time.Duration(intFromEnv(constants.ThundraLambdaTimeoutMargin,
		getDefaultTimeoutMargin())) * time.Millisecond
}

func boolFromEnv(key string, defaultValue bool) bool {
	env := os.Getenv(key)
	value, err := strconv.ParseBool(env)
	if err != nil {
		if env != "" {
			log.Printf("%v: %s is not a bool value", err, key)
		}
		return defaultValue
	}
	return value
}

func intFromEnv(key string, defaultValue int) int {
	t := os.Getenv(key)
	// environment variable is not set
	if t == "" {
		return defaultValue
	}

	i, err := strconv.Atoi(t)

	// environment variable is not set in the correct format
	if err != nil {
		log.Printf("%v: %s should be set with an integer\n", err, key)
		return defaultValue
	}
	return i
}

func determineAPIKey() string {
	apiKey := os.Getenv(constants.ThundraAPIKey)
	if apiKey == "" {
		log.Println("Error no APIKey in env variables")
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

func determineLogLevel() string {
	level := os.Getenv(constants.ThundraLogLogLevel)
	return strings.ToUpper(level)
}

func getDefaultTimeoutMargin() int {
	region := AwsLambdaRegion
	memory := AwsLambdaFunctionMemorySize
	timeoutMargin := 1000

	if region == "us-west-2" {
		timeoutMargin = 200
	} else if strings.HasPrefix(region, "us-west-") {
		timeoutMargin = 400
	} else if strings.HasPrefix(region, "us-") || strings.HasPrefix(region, "ca-") {
		timeoutMargin = 600
	} else if strings.HasPrefix(region, "ca-") {
		timeoutMargin = 800
	}

	normalizedTimeoutMargin := int((384.0 / float64(memory)) * float64(timeoutMargin))
	if normalizedTimeoutMargin > timeoutMargin {
		return normalizedTimeoutMargin
	}
	return timeoutMargin
}
