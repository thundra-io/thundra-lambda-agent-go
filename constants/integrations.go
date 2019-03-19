package constants

const HTTPMethodTag = "http.method"
const HTTPURLTag = "http.url"
const HTTPPathTag = "http.path"
const HTTPHostTag = "http.host"
const HTTPStatusTag = "http.status_code"
const HTTPQueryParamsTag = "http.query_params"

const HTTPClassName = "HTTP"
const HTTPDomainName = "API"

var DynamoDBRequestTypes = map[string]string {
	"BatchGetItem": "READ",
    "BatchWriteItem": "WRITE",
    "CreateTable": "WRITE",
    "CreateGlobalTable": "WRITE",
    "DeleteItem": "DELETE",
    "DeleteTable": "DELETE",
    "GetItem": "READ",
    "PutItem": "WRITE",
    "Query": "READ",
    "Scan": "READ",
    "UpdateItem": "WRITE",
}

var SpanTags = map[string]string {
	"OPERATION_TYPE": "operation.type",
    "TRIGGER_DOMAIN_NAME": "trigger.domainName",
    "TRIGGER_CLASS_NAME": "trigger.className",
    "TRIGGER_OPERATION_NAMES": "trigger.operationNames",
    "TOPOLOGY_VERTEX": "topology.vertex",
}

var DBTags = map[string]string {
	"DB_STATEMENT": "db.statement",
	"DB_STATEMENT_TYPE": "db.statement.type",
    "DB_INSTANCE": "db.instance",
    "DB_TYPE": "db.type",
    "DB_HOST": "db.host",
    "DB_PORT": "db.port",
    "DB_USER": "db.user",
}

var AwsDynamoDBTags = map[string]string {
    "TABLE_NAME": "aws.dynamodb.table.name",
    "REQUEST_THROTTLED": "aws.dynamodb.request.throttled",
}

var AwsSDKTags = map[string]string {
    "SERVICE_NAME": "aws.service.name",
    "REQUEST_NAME": "aws.request.name",
    "HOST": "host",
}