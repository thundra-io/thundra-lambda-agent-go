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

var EsTags = map[string]string{
	"ES_URI":    "elasticsearch.uri",
	"ES_METHOD": "elasticsearch.method",
	"ES_PARAMS": "elasticsearch.params",
	"ES_BODY":   "elasticsearch.body",
	"ES_HOSTS":  "elasticsearch.hosts",
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
	"MONGODB":       "MONGODB",
}

var MongoDBTags = map[string]string{
	"MONGODB_COMMAND":      "mongodb.command",
	"MONGODB_COMMAND_NAME": "mongodb.command.name",
	"MONGODB_COLLECTION":   "mongodb.collection.name",
}

var MongoDBCommandTypes = map[string]string{
	// Aggregate Commands
	"AGGREGATE": "READ",
	"COUNT":     "READ",
	"DISTINCT":  "READ",
	"GROUP":     "READ",
	"MAPREDUCE": "READ",

	// Geospatial Commands
	"GEONEAR":   "READ",
	"GEOSEARCH": "READ",

	// Query and Write Operation Commands
	"DELETE":                 "DELETE",
	"EVAL":                   "EXECUTE",
	"FIND":                   "READ",
	"FINDANDMODIFY":          "WRITE",
	"GETLASTERROR":           "READ",
	"GETMORE":                "READ",
	"GETPREVERROR":           "READ",
	"INSERT":                 "WRITE",
	"PARALLELCOLLECTIONSCAN": "READ",
	"RESETERROR":             "WRITE",
	"UPDATE":                 "WRITE",

	// Query Plan Cache Commands
	"PLANCACHECLEAR":           "DELETE",
	"PLANCACHECLEARFILTERS":    "DELETE",
	"PLANCACHELISTFILTERS":     "READ",
	"PLANCACHELISTPLANS":       "READ",
	"PLANCACHELISTQUERYSHAPES": "READ",
	"PLANCACHESETFILTER":       "WRITE",

	// Authentication Commands
	"AUTHENTICATE": "EXECUTE",
	"LOGOUT":       "EXECUTE",

	// User Management Commands
	"CREATEUSER":               "WRITE",
	"DROPALLUSERSFROMDATABASE": "DELETE",
	"DROPUSER":                 "DELETE",
	"GRANROLESTOUSER":          "WRITE",
	"REVOKEROLESFROMUSER":      "WRITE",
	"UPDATEUSER":               "WRITE",
	"USERSINFO":                "READ",

	// Role Management Commands
	"CREATEROLE":               "WRITE",
	"DROPROLE":                 "DELETE",
	"DROPALLROLESFROMDATABASE": "DELETE",
	"GRANTPRIVILEGESTOROLE":    "WRITE",
	"GRANTROLESTOROLE":         "WRITE",
	"INVALIDATEUSERCACHE":      "DELETE",
	"REVOKEPRIVILEGESFROMROLE": "WRITE",
	"REVOKEROLESFROMROLE":      "WRITE",
	"ROLESINFO":                "READ",
	"UPDATEROLE":               "WRITE",

	// Replication Commands
	"ISMASTER":                   "READ",
	"REPLSETABORTPRIMARYCATCHUP": "EXECUTE",
	"REPLSETFREEZE":              "EXECUTE",
	"REPLSETGETCONFIG":           "READ",
	"REPLSETGETSTATUS":           "READ",
	"REPLSETINITIATE":            "EXECUTE",
	"REPLSETMAINTENANCE":         "EXECUTE",
	"REPLSETRECONFIG":            "EXECUTE",
	"REPLSETRESIZEOPLOG":         "EXECUTE",
	"REPLSETSTEPDOWN":            "EXECUTE",
	"REPLSETSYNCFROM":            "EXECUTE",

	// Sharding Commands
	"ADDSHARD":            "EXECUTE",
	"ADDSHARDTOZONE":      "EXECUTE",
	"BALANCERSTART":       "EXECUTE",
	"BALANCERSTATUS":      "READ",
	"BALANCERSTOP":        "EXECUTE",
	"CLEANUPORPHANED":     "EXECUTE",
	"ENABLESHARDING":      "EXECUTE",
	"FLUSHROUTERCONFIG":   "EXECUTE",
	"ISDBGRID":            "READ",
	"LISTSHARDS":          "READ",
	"MOVEPRIMARY":         "EXECUTE",
	"MERGECHUNKS":         "EXECUTE",
	"REMOVESHARD":         "EXECUTE",
	"REMOVESHARDFROMZONE": "EXECUTE",
	"SHARDCOLLECTION":     "EXECUTE",
	"SHARDINGSTATE":       "READ",
	"SPLIT":               "EXECUTE",
	"UPDATEZONEKEYRANGE":  "EXECUTE",

	// Session Commands
	"ABORTTRANSACTION":         "EXECUTE",
	"COMMITTRANSACTION":        "EXECUTE",
	"ENDSESSIONS":              "EXECUTE",
	"KILLALLSESSIONS":          "EXECUTE",
	"KILLALLSESSIONSBYPATTERN": "EXECUTE",
	"KILLSESSIONS":             "EXECUTE",
	"REFRESHSESSIONS":          "EXECUTE",
	"STARTSESSION":             "EXECUTE",

	// Administration Commands
	"CLONE":                          "EXECUTE",
	"CLONECOLLECTION":                "EXECUTE",
	"CLONECOLLECTIONASCAPPED":        "EXECUTE",
	"COLLMOD":                        "WRITE",
	"COMPACT":                        "EXECUTE",
	"CONVERTTOCAPPED":                "EXECUTE",
	"COPYDB":                         "EXECUTE",
	"CREATE":                         "WRITE",
	"CREATEINDEXES":                  "WRITE",
	"CURRENTOP":                      "READ",
	"DROP":                           "DELETE",
	"DROPDATABASE":                   "DELETE",
	"DROPINDEXES":                    "DELETE",
	"FILEMD5":                        "READ",
	"FSYNC":                          "EXECUTE",
	"FSYNCUNLOCK":                    "EXECUTE",
	"GETPARAMETER":                   "READ",
	"KILLCURSORS":                    "EXECUTE",
	"KILLOP":                         "EXECUTE",
	"LISTCOLLECTIONS":                "READ",
	"LISTDATABASES":                  "READ",
	"LISTINDEXES":                    "READ",
	"LOGROTATE":                      "EXECUTE",
	"REINDEX":                        "WRITE",
	"RENAMECOLLECTION":               "WRITE",
	"REPAIRDATABASE":                 "EXECUTE",
	"SETFEATURECOMPATIBILITYVERSION": "WRITE",
	"SETPARAMETER":                   "WRITE",
	"SHUTDOWN":                       "EXECUTE",
	"TOUCH":                          "EXECUTE",

	// Diagnostic Commands
	"BUILDINFO":          "READ",
	"COLLSTATS":          "READ",
	"CONNPOOLSTATS":      "READ",
	"CONNECTIONSTATUS":   "READ",
	"CURSORINFO":         "READ",
	"DBHASH":             "READ",
	"DBSTATS":            "READ",
	"DIAGLOGGING":        "READ",
	"EXPLAIN":            "READ",
	"FEATURES":           "READ",
	"GETCMDLINEOPTS":     "READ",
	"GETLOG":             "READ",
	"HOSTINFO":           "READ",
	"LISTCOMMANDS":       "READ",
	"PROFILE":            "READ",
	"SERVERSTATUS":       "READ",
	"SHARDCONNPOOLSTATS": "READ",
	"TOP":                "READ",

	// Free Monitoring Commands
	"SETFREEMONITORING": "EXECUTE",

	// Auditing Commands
	"LOGAPPLICATIONMESSAGE": "EXECUTE",
}
