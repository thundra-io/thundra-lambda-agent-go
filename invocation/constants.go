package invocation

const invocationType = "Invocation"
const pluginOrder = 10
const defaultErrorCode = "-1"

const lambdaTriggerOperationName = "x-thundra-trigger-operation-name"

var classNames = map[string]string{
	"AWSSERVICE":    "AWSService",
	"DYNAMODB":      "AWS-DynamoDB",
	"SQS":           "AWS-SQS",
	"SNS":           "AWS-SNS",
	"KINESIS":       "AWS-Kinesis",
	"FIREHOSE":      "AWS-Firehose",
	"S3":            "AWS-S3",
	"LAMBDA":        "AWS-Lambda",
	"RDB":           "RDB",
	"REDIS":         "Redis",
	"HTTP":          "HTTP",
	"MYSQL":         "MYSQL",
	"POSTGRESQL":    "POSTGRESQL",
	"ELASTICSEARCH": "ELASTICSEARCH",
	"CLOUDFRONT":    "AWS-CloudFront",
	"APIGATEWAY":    "AWS-APIGateway",
	"CLOUDWATCHLOG": "AWS-CloudWatch-Log",
	"SCHEDULE":      "AWS-CloudWatch-Schedule",
}

var domainNames = map[string]string{
	"AWS":       "AWS",
	"DB":        "DB",
	"MESSAGING": "Messaging",
	"STREAM":    "Stream",
	"STORAGE":   "Storage",
	"API":       "API",
	"CACHE":     "Cache",
	"CDN":       "CDN",
	"LOG":       "LOG",
	"SCHEDULE":  "Schedule",
}

var spanTags = map[string]string{
	"OPERATION_TYPE":          "operation.type",
	"DB_INSTANCE":             "db.instance",
	"DB_TYPE":                 "db.type",
	"DB_HOST":                 "db.host",
	"TRIGGER_DOMAIN_NAME":     "trigger.domainName",
	"TRIGGER_CLASS_NAME":      "trigger.className",
	"TRIGGER_OPERATION_NAMES": "trigger.operationNames",
	"TOPOLOGY_VERTEX":         "topology.vertex",
	"DB_STATEMENT":            "db.statement",
	"DB_STATEMENT_TYPE":       "db.statement.type:",
}
