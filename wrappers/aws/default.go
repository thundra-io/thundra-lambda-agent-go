package thundraaws

import (
	"strings"

	"github.com/aws/aws-sdk-go/aws/request"
	"github.com/thundra-io/thundra-lambda-agent-go/constants"
	"github.com/thundra-io/thundra-lambda-agent-go/tracer"
)

type defaultAWSIntegration struct {
	ServiceName string
}

func newAWSServiceIntegration(serviceName string) *defaultAWSIntegration {
	return &defaultAWSIntegration{
		ServiceName: serviceName,
	}
}

func (i *defaultAWSIntegration) getOperationName(r *request.Request) string {
	return constants.AWSServiceRequest
}

func (i *defaultAWSIntegration) beforeCall(r *request.Request, span *tracer.RawSpan) {
	span.ClassName = constants.ClassNames["AWSSERVICE"]
	span.DomainName = constants.DomainNames["API"]

	tags := map[string]interface{}{
		constants.AwsSDKTags["REQUEST_NAME"]: strings.ToLower(i.ServiceName),
	}
	span.Tags = tags
}

func (i *defaultAWSIntegration) afterCall(r *request.Request, span *tracer.RawSpan) {

}
