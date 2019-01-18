package plugin

// const TimeFormat = "2006-01-02 15:04:05.000 -0700"
const DefaultProfile = "default"

//Thundra Environment
const ThundraApplicationProfile = "thundra_applicationProfile"
const AgentVersion = "2.0.0"
const DataModelVersion = "2.0"
const ApplicationDomainName = "API"
const ApplicationClassName = "AWS-Lambda"
const DefaultClassName = "Method"
const DefaultDomainName = ""
const ApplicationRuntime = "go"
const ApplicationRuntimeVersion = "1.x"
const ThundraDisableTrace = "thundra_agent_lambda_trace_disable"
const ThundraDisableMetric = "thundra_agent_lambda_metric_disable"
const ThundraDisableLog = "thundra_agent_lambda_log_disable"

//AWS
const AwsFunctionPlatform = "AWS Lambda"
const AwsDefaultRegion = "AWS_DEFAULT_REGION"
const thundraLambdaDebugEnable = "thundra_lambda_debug_enable"

//Pre-defined AWS Lambda Tags
const AwsLambdaInvocationRequestId = "aws.lambda.invocation.request_id"
const AwsLambdaInvocationRequest = "aws.lambda.invocation.request"
const AwsLambdaInvocationResponse = "aws.lambda.invocation.response"
const AwsLambdaARN = "aws.lambda.arn"
const AwsLambdaInvocationColdStart = "aws.lambda.invocation.coldstart"
const AwsLambdaInvocationTimeout = "aws.lambda.invocation.timeout"
const AwsLambdaLogGroupName = "aws.lambda.log_group_name"
const AwsLambdaLogStreamName = "aws.lambda.log_stream_name"
const AwsLambdaMemoryLimit = "aws.lambda.memory_limit"
const AwsLambdaName = "aws.lambda.name"
const AwsRegion = "aws.region"

const AwsError = "error"
const AwsErrorKind = "error.kind"
const AwsErrorMessage = "error.message"
const AwsErrorStack = "error.stack"
