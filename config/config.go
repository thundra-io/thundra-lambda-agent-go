package config

import (
	"fmt"
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
}

func isThundraDisabled() bool {
	env := os.Getenv(constants.ThundraLambdaDisable)
	disabled, err := strconv.ParseBool(env)
	if err != nil {
		if env != "" {
			log.Println(err, " thundra_lambda_disable is not a bool value. Thundra is enabled by default.")
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
			log.Println(err, constants.ThundraDisableTrace+" is not a bool value. Trace plugin is enabled by default.")
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
			log.Println(err, constants.ThundraDisableMetric+" is not a bool value. Metric plugin is enabled by default.")
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
			log.Println(err, constants.ThundraDisableLog+" is not a bool value. Log plugin is enabled by default.")
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

func determineRestCompositeBatchSize() int {
	t := os.Getenv(constants.ThundraLambdaReportRestCompositeBatchSize)
	// environment variable is not set
	if t == "" {
		return constants.ThundraLambdaReportRestCompositeBatchSizeDefault
	}

	i, err := strconv.Atoi(t)

	// environment variable is not set in the correct format
	if err != nil {
		fmt.Printf("%v: %s should be set with an integer\n", err, constants.ThundraLambdaReportRestCompositeBatchSize)
		return constants.ThundraLambdaReportRestCompositeBatchSizeDefault
	}

	return i
}

func determineCloudWatchCompositeBatchSize() int {
	t := os.Getenv(constants.ThundraLambdaReportCloudwatchCompositeBatchSize)
	// environment variable is not set
	if t == "" {
		return constants.ThundraLambdaReportCloudwatchCompositeBatchSizeDefault
	}

	i, err := strconv.Atoi(t)

	// environment variable is not set in the correct format
	if err != nil {
		fmt.Printf("%v: %s should be set with an integer\n", err, constants.ThundraLambdaReportCloudwatchCompositeBatchSize)
		return constants.ThundraLambdaReportCloudwatchCompositeBatchSizeDefault
	}

	return i
}

func determineWarmup() bool {
	w := os.Getenv(constants.ThundraLambdaWarmupWarmupAware)
	b, err := strconv.ParseBool(w)
	if err != nil {
		if w != "" {
			log.Println(err, " thundra_lambda_warmup_warmupAware should be set with a boolean.")
		}
		return false
	}
	return b
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
	env := os.Getenv(constants.ThundraDisableTraceRequest)
	disabled, err := strconv.ParseBool(env)
	if err != nil {
		if env != "" {
			log.Println(err, constants.ThundraDisableTraceRequest+"is not a bool value. Trace request is not disabled.")
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
			log.Println(err, constants.ThundraDisableTraceResponse+" is not a bool value. Trace response is not disabled.")
		}
		return false
	}
	return disabled
}

func determineLogLevel() string {
	level := os.Getenv(constants.ThundraLogLogLevel)
	return strings.ToUpper(level)
}

func isAwsIntegrationDisabled() bool {
	env := os.Getenv(constants.ThundraDisableAwsIntegration)
	disabled, err := strconv.ParseBool(env)
	if err != nil {
		if env != "" {
			log.Println(err, constants.ThundraDisableAwsIntegration+" is not a bool value.")
		}
		return false
	}
	return disabled
}

func isDynamoDBStatementsMasked() bool {
	env := os.Getenv(constants.ThundraMaskDynamoDBStatement)
	masked, err := strconv.ParseBool(env)
	if err != nil {
		if env != "" {
			log.Println(err, constants.ThundraMaskDynamoDBStatement+" is not a bool value.")
		}
		return false
	}
	return masked
}

func isRDBStatementsMasked() bool {
	env := os.Getenv(constants.ThundraMaskRDBStatement)
	masked, err := strconv.ParseBool(env)
	if err != nil {
		if env != "" {
			log.Println(err, constants.ThundraMaskRDBStatement+" is not a bool value.")
		}
		return false
	}
	return masked
}

func isEsBodyMasked() bool {
	env := os.Getenv(constants.ThundraMaskEsBody)
	masked, err := strconv.ParseBool(env)
	if err != nil {
		if env != "" {
			log.Println(err, constants.ThundraMaskEsBody+" is not a bool value.")
		}
		return false
	}
	return masked
}

func isRedisCommandMasked() bool {
	env := os.Getenv(constants.ThundraMaskRedisCommand)
	masked, err := strconv.ParseBool(env)
	if err != nil {
		if env != "" {
			log.Println(err, constants.ThundraMaskRedisCommand+" is not a bool value.")
		}
		return false
	}
	return masked
}

func isMongoDBCommandMasked() bool {
	env := os.Getenv(constants.ThundraMaskMongoDBCommand)
	masked, err := strconv.ParseBool(env)
	if err != nil {
		if env != "" {
			log.Println(err, constants.ThundraMaskMongoDBCommand+" is not a bool value.")
		}
		return false
	}
	return masked
}

func isTraceKinesisRequestEnabled() bool {
	env := os.Getenv(constants.ThundraLambdaTraceKinesisRequestEnable)
	enabled, err := strconv.ParseBool(env)
	if err != nil {
		if env != "" {
			log.Println(err, constants.ThundraLambdaTraceKinesisRequestEnable+" is not a bool value.")
		}
		return false
	}
	return enabled
}

func isTraceFirehoseRequestEnabled() bool {
	env := os.Getenv(constants.ThundraLambdaTraceFirehoseRequestEnable)
	enabled, err := strconv.ParseBool(env)
	if err != nil {
		if env != "" {
			log.Println(err, constants.ThundraLambdaTraceFirehoseRequestEnable+" is not a bool value.")
		}
		return false
	}
	return enabled
}

func isTraceCloudwatchlogRequestEnabled() bool {
	env := os.Getenv(constants.ThundraLambdaTraceCloudwatchlogRequestEnable)
	enabled, err := strconv.ParseBool(env)
	if err != nil {
		if env != "" {
			log.Println(err, constants.ThundraLambdaTraceCloudwatchlogRequestEnable+" is not a bool value.")
		}
		return false
	}
	return enabled
}

func isDynamoDBTraceInjectionEnabled() bool {
	env := os.Getenv(constants.EnableDynamoDbTraceInjection)
	enabled, err := strconv.ParseBool(env)
	if err != nil {
		if env != "" {
			log.Println(err, constants.EnableDynamoDbTraceInjection+" is not a bool value.")
		}
		return false
	}
	return enabled
}

func isLambdaTraceInjectionDisabled() bool {
	env := os.Getenv(constants.DisableLambdaTraceInjection)
	disabled, err := strconv.ParseBool(env)
	if err != nil {
		if env != "" {
			log.Println(err, constants.DisableLambdaTraceInjection+" is not a bool value.")
		}
		return false
	}
	return disabled
}

func isCloudwatchlogCompositeDataEnabled() bool {
	env := os.Getenv(constants.ThundraLambdaReportCloudwatchCompositeEnable)
	enabled, err := strconv.ParseBool(env)
	if err != nil {
		if env != "" {
			log.Println(err, constants.ThundraLambdaReportCloudwatchCompositeEnable+" is not a bool value.")
		}
		return true
	}
	return enabled
}

func isRestCompositeDataEnabled() bool {
	env := os.Getenv(constants.ThundraLambdaReportRestCompositeEnable)
	enabled, err := strconv.ParseBool(env)
	if err != nil {
		if env != "" {
			log.Println(err, constants.ThundraLambdaReportRestCompositeEnable+" is not a bool value.")
		}
		return true
	}
	return enabled
}

func isReportCloudwatchEnabled() bool {
	env := os.Getenv(constants.ThundraLambdaReportCloudwatchEnable)
	enabled, err := strconv.ParseBool(env)
	if err != nil {
		if env != "" {
			log.Println(err, constants.ThundraLambdaReportCloudwatchEnable+" is not a bool value.")
		}
		return false
	}
	return enabled
}
