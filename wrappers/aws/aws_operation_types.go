package thundraaws

import "regexp"

var awsOperationTypesExclusions = map[string]map[string]string{
	"AWS-Lambda": {
		"ListTags":          "READ",
		"TagResource":       "WRITE",
		"UntagResource":     "WRITE",
		"EnableReplication": "PERMISSION",
	},
	"AWS-S3": {
		"HeadBucket":                       "LIST",
		"ListBucket":                       "LIST",
		"ListBucketByTags":                 "READ",
		"ListBucketMultipartUploads":       "READ",
		"ListBucketVersions":               "READ",
		"ListJobs":                         "READ",
		"ListMultipartUploadParts":         "READ",
		"GetBucketTagging":                 "READ",
		"GetObjectVersionTagging":          "READ",
		"GetObjectTagging":                 "READ",
		"GetObject":                        "READ",
		"GetObjectAcl":                     "READ",
		"SelectObjectContent":              "READ",
		"DeleteObjectTagging":              "TAGGING",
		"DeleteObjectVersionTagging":       "TAGGING",
		"PutBucketTagging":                 "TAGGING",
		"PutObjectTagging":                 "TAGGING",
		"PutObjectVersionTagging":          "TAGGING",
		"AbortMultipartUpload":             "WRITE",
		"ReplicateDelete":                  "WRITE",
		"ReplicateObject":                  "WRITE",
		"RestoreObject":                    "WRITE",
		"GetBucketObjectLockConfiguration": "WRITE",
		"GetObjectLegalHold":               "WRITE",
		"GetObjectRetention":               "WRITE",
		"DeleteBucket":                     "WRITE",
		"CreateBucket":                     "WRITE",
		"DeleteObject":                     "WRITE",
		"PutObject":                        "WRITE",
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
		"GetRecords":                   "READ",
		"UntagResource":                "TAGGING",
		"ConfirmSubscription":          "WRITE",
		"OptInPhoneNumber":             "WRITE",
		"Subscribe":                    "WRITE",
		"Unsubscribe":                  "WRITE",
		"Publish":                      "WRITE",
		"PutRecords":                   "WRITE",
		"PutRecord":                    "WRITE",
	},
	"AWS-Athena": {
		"BatchGetNamedQuery":     "READ",
		"BatchGetQueryExecution": "READ",
		"GetNamedQuery":          "READ",
		"GetQueryExecution":      "READ",
		"GetQueryResults":        "READ",
		"GetWorkGroup":           "READ",
		"ListTagsForResource":    "READ",
		"CreateWorkGroup":        "WRITE",
		"CreateNamedQuery":       "WRITE",
		"DeleteNamedQuery":       "WRITE",
		"DeleteWorkGroup":        "WRITE",
		"UntagResource":          "TAGGING",
		"TagResource":            "TAGGING",
		"CancelQueryExecution":   "WRITE",
		"RunQuery":               "WRITE",
		"StartQueryExecution":    "WRITE",
		"StopQueryExecution":     "WRITE",
		"UpdateWorkGroup":        "WRITE",
		"ListNamedQueries":       "LIST",
		"ListQueryExecutions":    "LIST",
		"ListWorkGroups":         "LIST",
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
		"PutRecordBatch":                "WRITE",
		"PutRecord":                     "WRITE",
	},
	"AWS-SQS": {
		"ListDeadLetterSourceQueues": "READ",
		"ListQueueTags":              "READ",
		"ReceiveMessage":             "READ",
		"TagQueue":                   "TAGGING",
		"UntagQueue":                 "TAGGING",
		"PurgeQueue":                 "WRITE",
		"SetQueueAttributes":         "WRITE",
		"SendMessage":                "WRITE",
		"SendMessageBatch":           "WRITE",
		"DeleteMessage":              "WRITE",
		"DeleteMessageBatch":         "WRITE",
	},
	"AWS-DynamoDB": {
		"BatchGetItem":                      "READ",
		"ConditionCheckItem":                "READ",
		"ListStreams":                       "READ",
		"ListTagsOfResource":                "READ",
		"Query":                             "READ",
		"Scan":                              "READ",
		"GetItem":                           "READ",
		"TagResource":                       "TAGGING",
		"UntagResource":                     "TAGGING",
		"BatchWriteItem":                    "WRITE",
		"PurchaseReservedCapacityOfferings": "WRITE",
		"RestoreTableFromBackup":            "WRITE",
		"RestoreTableToPointInTime":         "WRITE",
		"CreateTable":                       "WRITE",
		"CreateGlobalTable":                 "WRITE",
		"DeleteItem":                        "WRITE",
		"DeleteTable":                       "WRITE",
		"UpdateItem":                        "WRITE",
		"PutItem":                           "WRITE",
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

var compiledTypes = map[string]regexp.Regexp{}