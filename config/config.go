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

func init() {
	ThundraDisabled = isThundraDisabled()
	TraceDisabled = isTraceDisabled()
	MetricDisabled = isMetricDisabled()
	LogDisabled = isLogDisabled()
	AwsIntegrationDisabled = isAwsIntegrationDisabled()
	TraceRequestDisabled = isTraceRequestDisabled()
	TraceResponseDisabled = isTraceResponseDisabled()
	DebugEnabled = isThundraDebugEnabled()
	TimeoutMargin = determineTimeoutMargin()
	WarmupEnabled = determineWarmup()
	APIKey = determineAPIKey()
	LogLevel = determineLogLevel()
	TrustAllCertificates = trustAllCertificates()
	MaskDynamoDBStatement = isDynamoDBStatementsMasked()
	MaskRDBStatement = isRDBStatementsMasked()
	MaskEsBody = isEsBodyMasked()
	MaskRedisCommand = isRedisCommandMasked()
	TraceKinesisRequestEnabled = isTraceKinesisRequestEnabled()
	TraceFirehoseRequestEnabled = isTraceFirehoseRequestEnabled()
	TraceCloudwatchlogRequestEnabled = isTraceCloudwatchlogRequestEnabled()
	DynamoDBTraceInjectionEnabled = isDynamoDBTraceInjectionEnabled()
	LambdaTraceInjectionDisabled = isLambdaTraceInjectionDisabled()
	ReportCloudwatchCompositeBatchSize = determineCloudWatchCompositeBatchSize()
	ReportRestCompositeBatchSize = determineRestCompositeBatchSize()
	ReportRestCompositeDataEnabled = isRestCompositeDataEnabled()
	ReportCloudwatchCompositeDataEnabled = isCloudwatchlogCompositeDataEnabled()
	ReportCloudwatchEnabled = isReportCloudwatchEnabled()
	MaskMongoDBCommand = isMongoDBCommandMasked()
	SamplingCountFrequency = determineSamplingCountFreq()
	SamplingTimeFrequency = determineSamplingTimeFreq()
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

func isThundraDisabled() bool {
	return boolFromEnv(constants.ThundraLambdaDisable, false)
}

func isTraceDisabled() bool {
	return boolFromEnv(constants.ThundraDisableTrace, false)
}

func isMetricDisabled() bool {
	return boolFromEnv(constants.ThundraDisableMetric, false)
}

func isLogDisabled() bool {
	return boolFromEnv(constants.ThundraDisableLog, false)
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
		log.Printf("%v: %s should be set with an integer\n", err, constants.ThundraLambdaTimeoutMargin)
		return time.Duration(constants.DefaultTimeoutMargin) * time.Millisecond
	}

	return time.Duration(i) * time.Millisecond
}

func determineRestCompositeBatchSize() int {
	return intFromEnv(constants.ThundraLambdaReportRestCompositeBatchSize, constants.ThundraLambdaReportRestCompositeBatchSizeDefault)
}

func determineCloudWatchCompositeBatchSize() int {
	return intFromEnv(constants.ThundraLambdaReportCloudwatchCompositeBatchSize, constants.ThundraLambdaReportCloudwatchCompositeBatchSizeDefault)
}

func determineWarmup() bool {
	return boolFromEnv(constants.ThundraLambdaWarmupWarmupAware, false)
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

func isTraceRequestDisabled() bool {
	return boolFromEnv(constants.ThundraDisableTraceRequest, false)
}

func isTraceResponseDisabled() bool {
	return boolFromEnv(constants.ThundraDisableTraceResponse, false)
}

func determineLogLevel() string {
	level := os.Getenv(constants.ThundraLogLogLevel)
	return strings.ToUpper(level)
}

func isAwsIntegrationDisabled() bool {
	return boolFromEnv(constants.ThundraDisableAwsIntegration, false)
}

func isDynamoDBStatementsMasked() bool {
	return boolFromEnv(constants.ThundraMaskDynamoDBStatement, false)
}

func isRDBStatementsMasked() bool {
	return boolFromEnv(constants.ThundraMaskRDBStatement, false)
}

func isEsBodyMasked() bool {
	return boolFromEnv(constants.ThundraMaskEsBody, false)
}

func isRedisCommandMasked() bool {
	return boolFromEnv(constants.ThundraMaskRedisCommand, false)
}

func isMongoDBCommandMasked() bool {
	return boolFromEnv(constants.ThundraMaskMongoDBCommand, false)
}

func isTraceKinesisRequestEnabled() bool {
	return boolFromEnv(constants.ThundraLambdaTraceKinesisRequestEnable, false)
}

func isTraceFirehoseRequestEnabled() bool {
	return boolFromEnv(constants.ThundraLambdaTraceFirehoseRequestEnable, false)
}

func isTraceCloudwatchlogRequestEnabled() bool {
	return boolFromEnv(constants.ThundraLambdaTraceCloudwatchlogRequestEnable, false)
}

func isDynamoDBTraceInjectionEnabled() bool {
	return boolFromEnv(constants.EnableDynamoDbTraceInjection, false)
}

func isLambdaTraceInjectionDisabled() bool {
	return boolFromEnv(constants.DisableLambdaTraceInjection, false)
}

func isCloudwatchlogCompositeDataEnabled() bool {
	return boolFromEnv(constants.ThundraLambdaReportCloudwatchCompositeEnable, true)
}

func isRestCompositeDataEnabled() bool {
	return boolFromEnv(constants.ThundraLambdaReportRestCompositeEnable, true)
}

func isReportCloudwatchEnabled() bool {
	return boolFromEnv(constants.ThundraLambdaReportCloudwatchEnable, false)
}

func determineSamplingCountFreq() int {
	return intFromEnv(constants.ThundraAgentMetricCountAwareSamplerCountFreq, -1)
}

func determineSamplingTimeFreq() int {
	return intFromEnv(constants.ThundraAgentMetricTimeAwareSamplerTimeFreq, -1)
}
