package trace

//Trace
const traceDataType = "AuditData"
const executionContext = "ExecutionContext"
const thundraLambdaTraceRequestDisable = "thundra_lambda_trace_request_disable"
const thundraLambdaTraceResponseDisable = "thundra_lambda_trace_response_disable"

//AuditInfo
const auditInfoOpenTimestamp = "openTimestamp"
const auditInfoCloseTimestamp = "closeTimestamp"
const auditInfoContextName = "contextName"
const auditInfoErrors = "errors"
const auditInfoId = "id"
const auditInfoThrownError = "thrownError"

//AuditInfoProperties
const auditInfoPropertiesRequest = "request"
const auditInfoPropertiesResponse = "response"
const auditInfoPropertiesColdStart = "coldStart"
const auditInfoPropertiesFunctionRegion = "functionRegion"
const auditInfoPropertiesFunctionMemoryLimit = "functionMemoryLimitInMB"

//AWS
const awsDefaultRegion = "AWS_DEFAULT_REGION"
