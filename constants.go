package thundra

const collectorUrl = "https://collector.thundra.io/api/monitor-datas"
const timeFormat = "2006-01-02 15:04:05.000 -0700"

//Thundra
const ThundraLambdaPublishCloudwatchEnable = "thundra_lambda_publish_cloudwatch_enable"
const ThundraApiKey = "thundra_apiKey"
const ThundraApplicationProfile = "thundra_applicationProfile"

//Trace
const dataFormatVersion = "1.0"
const dataType = "AuditData"
const executionContext = "ExecutionContext"
const applicationType = "go"
const defaultProfile = "default"

//AuditInfo
const auditInfoOpenTime = "openTime"
const audit_info_close_time = "closeTime"
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