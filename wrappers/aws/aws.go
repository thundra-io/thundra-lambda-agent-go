package thundraaws

import (
	opentracing "github.com/opentracing/opentracing-go"
	"github.com/thundra-io/thundra-lambda-agent-go/config"
	"github.com/thundra-io/thundra-lambda-agent-go/tracer"
	"github.com/thundra-io/thundra-lambda-agent-go/utils"
	"regexp"
	"strings"

	"github.com/aws/aws-sdk-go/aws/request"
	"github.com/aws/aws-sdk-go/aws/session"
)

// Wrap wraps the given session object and adds necessary
// handlers to create a span for the AWS call
func Wrap(s *session.Session) *session.Session {
	if !config.AwsIntegrationDisabled && s != nil {
		s.Handlers.Validate.PushFrontNamed(
			request.NamedHandler{
				Name: "github.com/thundra-io/thundra-lambda-agent-go/wrappers/aws/aws.go/validateHandler",
				Fn:   validateHandler,
			},
		)

		s.Handlers.Complete.PushFrontNamed(
			request.NamedHandler{
				Name: "github.com/thundra-io/thundra-lambda-agent-go/wrappers/aws/aws.go/completeHandler",
				Fn:   completeHandler,
			},
		)
	}
	return s
}

func validateHandler(r *request.Request) {
	serviceID := r.ClientInfo.ServiceID
	i, ok := integrations[serviceID]
	if !ok {
		i = newAWSServiceIntegration(serviceID)
	}
	span, ctxWithSpan := opentracing.StartSpanFromContext(r.Context(), i.getOperationName(r))
	r.SetContext(ctxWithSpan)
	rawSpan, ok := tracer.GetRaw(span)
	if !ok {
		return
	}
	i.beforeCall(r, rawSpan)
	tracer.OnSpanStarted(span)
}

func completeHandler(r *request.Request) {
	i, ok := integrations[r.ClientInfo.ServiceID]
	if !ok {
		return
	}
	span := opentracing.SpanFromContext(r.Context())
	if span == nil {
		return
	}
	rawSpan, ok := tracer.GetRaw(span)
	if !ok {
		return
	}
	i.afterCall(r, rawSpan)
	if r.Error != nil {
		utils.SetSpanError(span, r.Error)
	}
	span.Finish()
}

func getOperationType(operationName string, className string) string {
	operationName = strings.Title(operationName)
	if class, ok := awsOperationTypesExclusions[className]; ok {
		if exclusion, ok := class[operationName]; ok {
			return exclusion
		}
	}

	for pattern := range awsOperationTypesPatterns {
		if match, _ := regexp.MatchString(pattern, operationName); match {
			return awsOperationTypesPatterns[pattern]
		}
	}

	return ""
}
