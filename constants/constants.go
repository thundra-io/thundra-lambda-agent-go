package constants

const COLLECTOR_URL = "https://collector.thundra.io/api/monitor-datas"
const TIME_FORMAT = "2006-01-02 15:04:05.000 -0700"

//Thundra
const THUNDRA_LAMBDA_PUBLISH_CLOUDWATCH_ENABLE = "thundra_lambda_publish_cloudwatch_enable"
const THUNDRA_API_KEY = "thundra_apiKey"
const THUNDRA_APPLICATION_PROFILE = "thundra_applicationProfile"

//Trace
const DATA_FORMAT_VERSION = "1.0"
const DATA_TYPE = "AuditData"
const EXECUTION_CONTEXT = "ExecutionContext"
const APPLICATION_TYPE = "go"
const DEFAULT_PROFILE = "default"

//AuditInfo
const AUDIT_INFO_OPEN_TIME = "openTime"
const AUDIT_INFO_CLOSE_TIME = "closeTime"
const AUDIT_INFO_CONTEXT_NAME = "contextName"
const AUDIT_INFO_ERRORS = "errors"
const AUDIT_INFO_ID = "id"
const AUDIT_INFO_THROWN_ERROR = "thrownError"

//AuditInfoProperties
const AUDIT_INFO_PROPERTIES_REQUEST = "request"
const AUDIT_INFO_PROPERTIES_RESPONSE = "response"
const AUDIT_INFO_PROPERTIES_COLD_START = "coldStart"
const AUDIT_INFO_PROPERTIES_FUNCTION_REGION = "functionRegion"
const AUDIT_INFO_PROPERTIES_FUNCTION_MEMORY_LIMIT = "functionMemoryLimitInMB"

//AWS
const AWS_DEFAULT_REGION = "AWS_DEFAULT_REGION"