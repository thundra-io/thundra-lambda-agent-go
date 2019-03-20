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

var SQSRequestTypes = map[string]string {
    "ReceiveMessage": "READ",
    "SendMessage": "WRITE",
    "SendMessageBatch": "WRITE",
    "DeleteMessage": "DELETE",
    "DeleteMessageBatch": "DELETE",
}

var SNSRequestTypes = map[string]string {
    "Publish": "WRITE",
}

var KinesisRequestTypes = map[string]string {
    "GetRecords": "READ",
    "PutRecords": "WRITE",
    "PutRecord": "WRITE",
}

var FirehoseRequestTypes = map[string]string {
    "PutRecordBatch": "WRITE",
    "PutRecord": "WRITE",
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

var AwsSQSTags = map[string]string {
    "QUEUE_NAME": "aws.sqs.queue.name",
}

var AwsSNSTags = map[string]string {
    "TOPIC_NAME": "aws.sns.topic.name",
}

var AwsKinesisTags = map[string]string {
    "STREAM_NAME": "aws.kinesis.stream.name",
}

var AwsFirehoseTags = map[string]string {
    "STREAM_NAME": "aws.firehose.stream.name",
}

var AwsSDKTags = map[string]string {
    "SERVICE_NAME": "aws.service.name",
    "REQUEST_NAME": "aws.request.name",
    "HOST": "host",
}

var DomainNames = map[string]string {
    "AWS": "AWS",
    "DB": "DB",
    "MESSAGING": "Messaging",
    "STREAM": "Stream",
    "STORAGE": "Storage",
    "API": "API",
    "CACHE": "Cache",
}

var ClassNames = map[string]string {
    "AWSSERVICE": "AWSService",
    "DYNAMODB": "AWS-DynamoDB",
    "SQS": "AWS-SQS",
    "SNS": "AWS-SNS",
    "KINESIS": "AWS-Kinesis",
    "FIREHOSE": "AWS-Firehose",
    "S3": "AWS-S3",
    "LAMBDA": "AWS-Lambda",
    "RDB": "RDB",
    "REDIS": "Redis",
    "HTTP": "HTTP",
    "MYSQL": "MYSQL",
    "POSTGRESQL": "POSTGRESQL",
    "ELASTICSEARCH": "ELASTICSEARCH",
}