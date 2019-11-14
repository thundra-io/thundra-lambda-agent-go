package thundraaws

var awsOperationTypesExclusions = map[string]map[string]string{
	"AWS-Lambda": {
		"ListTags":          "READ",
		"TagResource":       "WRITE",
		"UntagResource":     "WRITE",
		"EnableReplication": "PERMISSION",
	},
	"AWS-S3": {
		"HeadBucket":                       "LIST",
		"ListBucketByTags":                 "READ",
		"ListBucketMultipartUploads":       "READ",
		"ListBucketVersions":               "READ",
		"ListJobs":                         "READ",
		"ListMultipartUploadParts":         "READ",
		"GetBucketTagging":                 "READ",
		"GetObjectVersionTagging":          "READ",
		"GetObjectTagging":                 "READ",
		"GetBucketObjectLockConfiguration": "WRITE",
		"GetObjectLegalHold":               "WRITE",
		"GetObjectRetention":               "WRITE",
		"DeleteObjectTagging":              "TAGGING",
		"DeleteObjectVersionTagging":       "TAGGING",
		"PutBucketTagging":                 "TAGGING",
		"PutObjectTagging":                 "TAGGING",
		"PutObjectVersionTagging":          "TAGGING",
		"AbortMultipartUpload":             "WRITE",
		"ReplicateDelete":                  "WRITE",
		"ReplicateObject":                  "WRITE",
		"RestoreObject":                    "WRITE",
		"DeleteBucketPolicy":               "PERMISSION",
		"ObjectOwnerOverrideToBucketOwner": "PERMISSION",
		"PutAccountPublicAccessBlock":      "PERMISSION",
		"PutBucketAcl":                     "PERMISSION",
		"PutBucketPolicy":                  "PERMISSION",
		"PutBucketPublicAccessBlock":       "PERMISSION",
		"PutObjectAcl":                     "PERMISSION",
		"PutObjectVersionAcl":              "PERMISSION",
	},
	"AWS-SNS": {
		"ListPhoneNumbersOptedOut":     "READ",
		"ListTagsForResource":          "READ",
		"CheckIfPhoneNumberIsOptedOut": "READ",
		"UntagResource":                "TAGGING",
		"ConfirmSubscription":          "WRITE",
		"OptInPhoneNumber":             "WRITE",
		"Subscribe":                    "WRITE",
		"Unsubscribe":                  "WRITE",
	},
	"AWS-Athena": {
		"BatchGetNamedQuery":     "READ",
		"BatchGetQueryExecution": "READ",
		"ListTagsForResource":    "LIST",
		"CreateWorkGroup":        "WRITE",
		"UntagResource":          "TAGGING",
		"TagResource":            "TAGGING",
		"CancelQueryExecution":   "WRITE",
		"RunQuery":               "WRITE",
		"StartQueryExecution":    "WRITE",
		"StopQueryExecution":     "WRITE",
	},
	"AWS-Kinesis": {
		"ListTagsForStream":             "READ",
		"SubscribeToShard":              "READ",
		"AddTagsToStream":               "TAGGING",
		"RemoveTagsFromStream":          "TAGGING",
		"DecreaseStreamRetentionPeriod": "WRITE",
		"DeregisterStreamConsumer":      "WRITE",
		"DisableEnhancedMonitoring":     "WRITE",
		"EnableEnhancedMonitoring":      "WRITE",
		"IncreaseStreamRetentionPeriod": "WRITE",
		"MergeShards":                   "WRITE",
		"RegisterStreamConsumer":        "WRITE",
		"SplitShard":                    "WRITE",
		"UpdateShardCount":              "WRITE",
	},
	"AWS-Firehose": {
		"DescribeDeliveryStream":        "LIST",
		"StartDeliveryStreamEncryption": "WRITE",
		"StopDeliveryStreamEncryption":  "WRITE",
		"TagDeliveryStream":             "WRITE",
		"UntagDeliveryStream":           "WRITE",
	},
	"AWS-SQS": {
		"ListDeadLetterSourceQueues": "READ",
		"ListQueueTags":              "READ",
		"ReceiveMessage":             "READ",
		"TagQueue":                   "TAGGING",
		"UntagQueue":                 "TAGGING",
		"PurgeQueue":                 "WRITE",
		"SetQueueAttributes":         "WRITE",
	},
	"AWS-DynamoDB": {
		"BatchGetItem":                      "READ",
		"ConditionCheckItem":                "READ",
		"ListStreams":                       "READ",
		"ListTagsOfResource":                "READ",
		"Query":                             "READ",
		"Scan":                              "READ",
		"TagResource":                       "TAGGING",
		"UntagResource":                     "TAGGING",
		"BatchWriteItem":                    "WRITE",
		"PurchaseReservedCapacityOfferings": "WRITE",
		"RestoreTableFromBackup":            "WRITE",
		"RestoreTableToPointInTime":         "WRITE",
	},
}

var awsOperationTypesPatterns = map[string]string{
	"^List.*$":       "LIST",
	"^Get.*$":        "READ",
	"^Create.*$":     "WRITE",
	"^Delete.*$":     "WRITE",
	"^Invoke.*$":     "WRITE",
	"^Publish.*$":    "WRITE",
	"^Put.*$":        "WRITE",
	"^Update.*$":     "WRITE",
	"^Describe.*$":   "READ",
	"^Change.*$":     "WRITE",
	"^Send.*$":       "WRITE",
	"^.*Permission$": "PERMISSION",
	"^.*Tagging$":    "TAGGING",
	"^.*Tags$":       "TAGGING",
	"^Set.*$":        "WRITE",
}