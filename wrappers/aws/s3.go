package thundraaws

import (
	"encoding/json"

	"github.com/thundra-io/thundra-lambda-agent-go/application"
	"github.com/thundra-io/thundra-lambda-agent-go/constants"

	"github.com/aws/aws-sdk-go/aws/request"
	"github.com/thundra-io/thundra-lambda-agent-go/tracer"
)

type s3Integration struct{}

func (i *s3Integration) getBucketName(r *request.Request) string {
	fields := struct {
		Bucket string `json:"Bucket"`
	}{}
	m, err := json.Marshal(r.Params)
	if err != nil {
		return ""
	}
	if err = json.Unmarshal(m, &fields); err != nil {
		return ""
	}
	if len(fields.Bucket) > 0 {
		return fields.Bucket
	}
	return ""
}

func (i *s3Integration) getKeyName(r *request.Request) string {
	fields := struct {
		Key string `json:"Key"`
	}{}
	m, err := json.Marshal(r.Params)
	if err != nil {
		return ""
	}
	if err = json.Unmarshal(m, &fields); err != nil {
		return ""
	}
	if len(fields.Key) > 0 {
		return fields.Key
	}
	return ""
}

func (i *s3Integration) getOperationName(r *request.Request) string {
	bucketName := i.getBucketName(r)
	if len(bucketName) > 0 {
		return bucketName
	}
	return constants.AWSServiceRequest
}

func (i *s3Integration) beforeCall(r *request.Request, span *tracer.RawSpan) {
	span.ClassName = constants.ClassNames["S3"]
	span.DomainName = constants.DomainNames["STORAGE"]

	operationName := r.Operation.Name
	operationType := constants.S3RequestTypes[operationName]

	tags := map[string]interface{}{
		constants.AwsS3Tags["BUCKET_NAME"]:            i.getBucketName(r),
		constants.AwsS3Tags["OBJECT_NAME"]:            i.getKeyName(r),
		constants.SpanTags["OPERATION_TYPE"]:          operationType,
		constants.AwsSDKTags["REQUEST_NAME"]:          operationName,
		constants.SpanTags["TOPOLOGY_VERTEX"]:         true,
		constants.SpanTags["TRIGGER_OPERATION_NAMES"]: []string{application.FunctionName},
		constants.SpanTags["TRIGGER_DOMAIN_NAME"]:     constants.AwsLambdaApplicationDomain,
		constants.SpanTags["TRIGGER_CLASS_NAME"]:      constants.AwsLambdaApplicationClass,
	}

	span.Tags = tags
}

func (i *s3Integration) afterCall(r *request.Request, span *tracer.RawSpan) {
	return
}

func init() {
	integrations["S3"] = &s3Integration{}
}
