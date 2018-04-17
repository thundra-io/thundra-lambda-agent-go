package trace

//Trace
const TraceDataType = "AuditData"
const executionContext = "ExecutionContext"

//AuditInfo
const auditInfoOpenTime = "openTime"
const auditInfoCloseTime = "closeTime"
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