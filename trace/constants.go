package trace

const timeFormat = "2006-01-02 15:04:05.000 -0700"

//Trace
const TraceDataType = "AuditData"
const executionContext = "ExecutionContext"
const applicationType = "go"
const defaultProfile = "default"

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

//Thundra Environment
const thundraApplicationProfile = "thundra_applicationProfile"

//AWS
const awsDefaultRegion = "AWS_DEFAULT_REGION"