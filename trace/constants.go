package trace

//Trace
const traceDataType = "AuditData"
const executionContext = "ExecutionContext"
const thundraLambdaHideRequest = "thundra_lambda_hide_request"
const thundraLambdaHideResponse = "thundra_lambda_hide_response"

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