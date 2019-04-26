package constants

const AWSServiceRequest = "AWSServiceRequest"

var HTTPTags = map[string]string{
	"METHOD":       "http.method",
	"URL":          "http.url",
	"PATH":         "http.path",
	"HOST":         "http.host",
	"STATUS":       "http.status_code",
	"QUERY_PARAMS": "http.query_params",
}

var DynamoDBRequestTypes = map[string]string{
	"BatchGetItem":      "READ",
	"BatchWriteItem":    "WRITE",
	"CreateTable":       "WRITE",
	"CreateGlobalTable": "WRITE",
	"DeleteItem":        "DELETE",
	"DeleteTable":       "DELETE",
	"GetItem":           "READ",
	"PutItem":           "WRITE",
	"Query":             "READ",
	"Scan":              "READ",
	"UpdateItem":        "WRITE",
}

var SQSRequestTypes = map[string]string{
	"ReceiveMessage":     "READ",
	"SendMessage":        "WRITE",
	"SendMessageBatch":   "WRITE",
	"DeleteMessage":      "DELETE",
	"DeleteMessageBatch": "DELETE",
}

var SNSRequestTypes = map[string]string{
	"Publish": "WRITE",
}

var KinesisRequestTypes = map[string]string{
	"GetRecords": "READ",
	"PutRecords": "WRITE",
	"PutRecord":  "WRITE",
}

var FirehoseRequestTypes = map[string]string{
	"PutRecordBatch": "WRITE",
	"PutRecord":      "WRITE",
}

var S3RequestTypes = map[string]string{
	"DeleteBucket":  "DELETE",
	"CreateBucket":  "WRITE",
	"copyObject":    "WRITE",
	"DeleteObject":  "DELETE",
	"deleteObjects": "DELETE",
	"GetObject":     "READ",
	"GetObjectAcl":  "READ",
	"ListBucket":    "READ",
	"PutObject":     "WRITE",
	"PutObjectAcl":  "WRITE",
}

var LambdaRequestTypes = map[string]string{
	"InvokeAsync": "CALL",
	"Invoke":      "CALL",
}

var SpanTags = map[string]string{
	"OPERATION_TYPE":          "operation.type",
	"TRIGGER_DOMAIN_NAME":     "trigger.domainName",
	"TRIGGER_CLASS_NAME":      "trigger.className",
	"TRIGGER_OPERATION_NAMES": "trigger.operationNames",
	"TOPOLOGY_VERTEX":         "topology.vertex",
	"TRACE_LINKS":             "trace.links",
}

var DBTags = map[string]string{
	"DB_STATEMENT":      "db.statement",
	"DB_STATEMENT_TYPE": "db.statement.type",
	"DB_INSTANCE":       "db.instance",
	"DB_TYPE":           "db.type",
	"DB_HOST":           "db.host",
	"DB_PORT":           "db.port",
	"DB_USER":           "db.user",
}

var AwsDynamoDBTags = map[string]string{
	"TABLE_NAME":        "aws.dynamodb.table.name",
	"REQUEST_THROTTLED": "aws.dynamodb.request.throttled",
}

var AwsSQSTags = map[string]string{
	"QUEUE_NAME": "aws.sqs.queue.name",
}

var AwsSNSTags = map[string]string{
	"TOPIC_NAME": "aws.sns.topic.name",
}

var AwsKinesisTags = map[string]string{
	"STREAM_NAME": "aws.kinesis.stream.name",
}

var AwsFirehoseTags = map[string]string{
	"STREAM_NAME": "aws.firehose.stream.name",
}

var AwsSDKTags = map[string]string{
	"SERVICE_NAME": "aws.service.name",
	"REQUEST_NAME": "aws.request.name",
	"HOST":         "host",
}

var AwsS3Tags = map[string]string{
	"BUCKET_NAME": "aws.s3.bucket.name",
	"OBJECT_NAME": "aws.s3.object.name",
}

var AwsLambdaTags = map[string]string{
	"FUNCTION_NAME":      "aws.lambda.name",
	"FUNCTION_QUALIFIER": "aws.lambda.qualifier",
	"INVOCATION_TYPE":    "aws.lambda.invocation.type",
	"INVOCATION_PAYLOAD": "aws.lambda.invocation.payload",
}

var RedisCommandTypes = map[string]string{
	"APPEND":            "WRITE",
	"BGREWRITEAOF":      "WRITE",
	"BGSAVE":            "WRITE",
	"BITCOUNT":          "READ",
	"BITFIELD":          "WRITE",
	"BITOP":             "WRITE",
	"BITPOS":            "READ",
	"BLPOP":             "DELETE",
	"BRPOP":             "DELETE",
	"BRPOPLPUSH":        "WRITE",
	"BZPOPMIN":          "DELETE",
	"BZPOPMAX":          "DELETE",
	"DBSIZE":            "READ",
	"DECR":              "WRITE",
	"DECRBY":            "WRITE",
	"DELETE":            "DELETE",
	"EVAL":              "EXECUTE",
	"EVALSHA":           "EXECUTE",
	"EXISTS":            "READ",
	"EXPIRE":            "WRITE",
	"EXPIREAT":          "WRITE",
	"FLUSHALL":          "DELETE",
	"FLUSHDB":           "DELETE",
	"GEOADD":            "WRITE",
	"GEOHASH":           "READ",
	"GEOPOS":            "READ",
	"GEODIST":           "READ",
	"GEORADIUS":         "READ",
	"GEORADIUSBYMEMBER": "READ",
	"GET":               "READ",
	"GETBIT":            "READ",
	"GETRANGE":          "READ",
	"GETSET":            "WRITE",
	"HDEL":              "DELETE",
	"HEXISTS":           "READ",
	"HGET":              "READ",
	"HGETALL":           "READ",
	"HINCRBY":           "WRITE",
	"HINCRBYFLOAT":      "WRITE",
	"HKEYS":             "READ",
	"HLEN":              "READ",
	"HMGET":             "READ",
	"HMSET":             "WRITE",
	"HSET":              "WRITE",
	"HSETNX":            "WRITE",
	"HSTRLEN":           "READ",
	"HVALS":             "READ",
	"INCR":              "WRITE",
	"INCRBY":            "WRITE",
	"INCRBYFLOAT":       "WRITE",
	"KEYS":              "READ",
	"LINDEX":            "READ",
	"LINSERT":           "WRITE",
	"LLEN":              "READ",
	"LPOP":              "DELETE",
	"LPUSH":             "WRITE",
	"LPUSHX":            "WRITE",
	"LRANGE":            "READ",
	"LREM":              "DELETE",
	"LSET":              "WRITE",
	"LTRIM":             "DELETE",
	"MGET":              "READ",
	"MSET":              "WRITE",
	"MSETNX":            "WRITE",
	"PERSIST":           "WRITE",
	"PEXPIRE":           "WRITE",
	"PEXPIREAT":         "WRITE",
	"PFADD":             "WRITE",
	"PFCOUNT":           "READ",
	"PFMERGE":           "WRITE",
	"PSETEX":            "WRITE",
	"PUBLISH":           "WRITE",
	"RPOP":              "DELETE",
	"RPOPLPUSH":         "WRITE",
	"RPUSH":             "WRITE",
	"RPUSHX":            "WRITE",
	"SADD":              "WRITE",
	"SCARD":             "READ",
	"SDIFFSTORE":        "WRITE",
	"SET":               "WRITE",
	"SETBIT":            "WRITE",
	"SETEX":             "WRITE",
	"SETNX":             "WRITE",
	"SETRANGE":          "WRITE",
	"SINTER":            "READ",
	"SINTERSTORE":       "WRITE",
	"SISMEMBER":         "READ",
	"SMEMBERS":          "READ",
	"SMOVE":             "WRITE",
	"SORT":              "WRITE",
	"SPOP":              "DELETE",
	"SRANDMEMBER":       "READ",
	"SREM":              "DELETE",
	"STRLEN":            "READ",
	"SUNION":            "READ",
	"SUNIONSTORE":       "WRITE",
	"ZADD":              "WRITE",
	"ZCARD":             "READ",
	"ZCOUNT":            "READ",
	"ZINCRBY":           "WRITE",
	"ZINTERSTORE":       "WRITE",
	"ZLEXCOUNT":         "READ",
	"ZPOPMAX":           "DELETE",
	"ZPOPMIN":           "DELETE",
	"ZRANGE":            "READ",
	"ZRANGEBYLEX":       "READ",
	"ZREVRANGEBYLEX":    "READ",
	"ZRANGEBYSCORE":     "READ",
	"ZRANK":             "READ",
	"ZREM":              "DELETE",
	"ZREMRANGEBYLEX":    "DELETE",
	"ZREMRANGEBYRANK":   "DELETE",
	"ZREMRANGEBYSCORE":  "DELETE",
	"ZREVRANGE":         "READ",
	"ZREVRANGEBYSCORE":  "READ",
	"ZREVRANK":          "READ",
	"ZSCORE":            "READ",
	"ZUNIONSTORE":       "WRITE",
	"SCAN":              "READ",
	"SSCAN":             "READ",
	"HSCAN":             "READ",
	"ZSCAN":             "READ",
	"XADD":              "WRITE",
	"XRANGE":            "READ",
	"XREVRANGE":         "READ",
	"XLEN":              "READ",
	"XREAD":             "READ",
	"XREADGROUP":        "READ",
	"XPENDING":          "READ",
}

var RedisTags = map[string]string{
	"REDIS_HOST":         "redis.host",
	"REDIS_PORT":         "redis.port",
	"REDIS_COMMAND":      "redis.command",
	"REDIS_COMMANDS":     "redis.commands",
	"REDIS_COMMAND_TYPE": "redis.command.type",
	"REDIS_COMMAND_ARGS": "redis.command.args",
}

var DomainNames = map[string]string{
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

var ClassNames = map[string]string{
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
