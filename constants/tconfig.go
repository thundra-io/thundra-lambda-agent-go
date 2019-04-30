package constants

const AgentVersion = "2.0.0"
const DataModelVersion = "2.0"

const DefaultProfile = "default"
const DefaultClassName = "Method"
const DefaultDomainName = ""

const ThundraLambdaDisable = "thundra_agent_lambda_disable"
const ThundraDisableTrace = "thundra_agent_lambda_trace_disable"
const ThundraDisableMetric = "thundra_agent_lambda_metric_disable"
const ThundraDisableLog = "thundra_agent_lambda_log_disable"
const ThundraDisableTraceRequest = "thundra_agent_lambda_trace_request_disable"
const ThundraDisableTraceResponse = "thundra_agent_lambda_trace_response_disable"
const ThundraApplicationStage = "thundra_agent_lambda_application_stage"
const ThundraLogLogLevel = "thundra_log_logLevel"

const ThundraDisableAwsIntegration = "thundra_agent_lambda_trace_integrations_aws_disable"

const DefaultTimeoutMargin = 200
const DefaultCollectorURL = "https://api.thundra.io/v1"
const MonitoringDataPath = "/monitoring-data"
const CompositeMonitoringDataPath = "/composite-monitoring-data"

const ThundraAPIKey = "thundra_apiKey"
const ThundraLambdaPublishCloudwatchEnable = "thundra_agent_lambda_publish_cloudwatch_enable"
const ThundraLambdaReportRestBaseURL = "thundra_agent_lambda_report_rest_baseUrl"
const ThundraLambdaWarmupWarmupAware = "thundra_agent_lambda_warmup_warmupAware"
const ThundraLambdaTimeoutMargin = "thundra_agent_lambda_timeout_margin"
const ThundraTrustAllCertificates = "thundra_agent_lambda_report_rest_trustAllCertificates"

const ApplicationIDProp = "thundra_agent_lambda_application_id"
const ApplicationDomainProp = "thundra_agent_lambda_application_domainName"
const ApplicationClassProp = "thundra_agent_lambda_application_className"
const ApplicationNameProp = "thundra_agent_lambda_application_name"
const ApplicationVersionProp = "thundra_agent_lambda_application_version"
const ApplicationStageProp = "thundra_agent_lambda_application_stage"
const ApplicationTagPrefixProp = "thundra_agent_lambda_application_tag_"

const ThundraMaskDynamoDBStatement = "thundra_agent_lambda_trace_integrations_aws_dynamodb_statement_mask"
const ThundraMaskRDBStatement = "thundra_agent_lambda_trace_integrations_rdb_statement_mask"
const ThundraMaskRedisCommand = "thundra_agent_lambda_trace_integrations_redis_command_mask"

const ThundraLambdaTraceKinesisRequestEnable = "thundra_agent_lambda_trace_kinesis_request_enable"
const ThundraLambdaTraceFirehoseRequestEnable = "thundra_agent_lambda_trace_firehose_request_enable"
const ThundraLambdaTraceCloudwatchlogRequestEnable = "thundra_agent_lambda_trace_cloudwatchlog_request_enable"
const ThundraMaskEsBody = "thundra_agent_lambda_trace_integrations_elasticsearch_body_mask"

const EnableDynamoDbTraceInjection = "thundra_agent_trace_integrations_dynamodb_trace_injection_enable"
const DisableLambdaTraceInjection = "thundra_agent_trace_integrations_aws_lambda_traceInjection_disable"

const MaxTracedHttpBodySize = 128 * 1024
const ThundraLambdaReportRestCompositeBatchSize = "thundra_agent_lambda_report_rest_composite_batchsize"
const ThundraLambdaReportCloudwatchCompositeBatchSize = "thundra_agent_lambda_report_cloudwatch_composite_batchsize"

const ThundraLambdaReportRestCompositeEnable = "thundra_agent_lambda_report_rest_composite_enable"
const ThundraLambdaReportCloudwatchCompositeEnable = "thundra_agent_lambda_report_cloudwatch_composite_enable"

const ThundraLambdaReportRestCompositeBatchSizeDefault = 100
const ThundraLambdaReportCloudwatchCompositeBatchSizeDefault = 10
