package trace

const TimeFormat = "2006-01-02 15:04:05.000 -0700"

//TODO make them private

//Trace
const TraceDataType = "AuditData"
const ExecutionContext = "ExecutionContext"
const ApplicationType = "go"
const DefaultProfile = "default"

//AuditInfo
const AuditInfoOpenTime = "openTime"
const AuditInfoCloseTime = "closeTime"
const AuditInfoContextName = "contextName"
const AuditInfoErrors = "errors"
const AuditInfoId = "id"
const AuditInfoThrownError = "thrownError"

//AuditInfoProperties
const AuditInfoPropertiesRequest = "request"
const AuditInfoPropertiesResponse = "response"
const AuditInfoPropertiesColdStart = "coldStart"
const AuditInfoPropertiesFunctionRegion = "functionRegion"
const AuditInfoPropertiesFunctionMemoryLimit = "functionMemoryLimitInMB"

//Thundra Environment
const ThundraApplicationProfile = "thundra_applicationProfile"

//AWS
const AwsDefaultRegion = "AWS_DEFAULT_REGION"