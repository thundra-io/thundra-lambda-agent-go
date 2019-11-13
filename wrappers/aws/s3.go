package thundraaws

import (
	"encoding/json"

	"github.com/thundra-io/thundra-lambda-agent-go/application"
	"github.com/thundra-io/thundra-lambda-agent-go/constants"

	"github.com/aws/aws-sdk-go/aws/request"
	"github.com/thundra-io/thundra-lambda-agent-go/tracer"
)

type s3Integration struct{}
type s3Params struct {
	Bucket string
	Key    string
}

func (i *s3Integration) getS3Info(r *request.Request) *s3Params {
	fields := &s3Params{}
	m, err := json.Marshal(r.Params)
	if err != nil {
		return &s3Params{}
	}
	if err = json.Unmarshal(m, &fields); err != nil {
		return &s3Params{}
	}
	return fields
}

func (i *s3Integration) getOperationName(r *request.Request) string {
	s3Info := i.getS3Info(r)
	if len(s3Info.Bucket) > 0 {
		return s3Info.Bucket
	}
	return constants.AWSServiceRequest
}

func (i *s3Integration) beforeCall(r *request.Request, span *tracer.RawSpan) {
	span.ClassName = constants.ClassNames["S3"]
	span.DomainName = constants.DomainNames["STORAGE"]

	operationName := r.Operation.Name
	operationType := getOperationType(operationName, constants.ClassNames["S3"])

	s3Info := i.getS3Info(r)

	tags := map[string]interface{}{
		constants.SpanTags["OPERATION_TYPE"]:          operationType,
		constants.AwsSDKTags["REQUEST_NAME"]:          operationName,
		constants.SpanTags["TOPOLOGY_VERTEX"]:         true,
		constants.SpanTags["TRIGGER_OPERATION_NAMES"]: []string{application.FunctionName},
		constants.SpanTags["TRIGGER_DOMAIN_NAME"]:     constants.AwsLambdaApplicationDomain,
		constants.SpanTags["TRIGGER_CLASS_NAME"]:      constants.AwsLambdaApplicationClass,
	}

	if s3Info.Bucket != "" {
		tags[constants.AwsS3Tags["BUCKET_NAME"]] = s3Info.Bucket
	}
	if s3Info.Key != "" {
		tags[constants.AwsS3Tags["OBJECT_NAME"]] = s3Info.Key
	}

	span.Tags = tags
}

func (i *s3Integration) afterCall(r *request.Request, span *tracer.RawSpan) {
	xAmzRequestID := ""
	if r.HTTPResponse != nil && r.HTTPResponse.Header != nil {
		xAmzRequestID = r.HTTPResponse.Header.Get("x-amz-request-id")
	}
	if xAmzRequestID != "" {
		span.Tags[constants.SpanTags["TRACE_LINKS"]] = []string{xAmzRequestID}
	}
}

func init() {
	integrations["S3"] = &s3Integration{}
}
