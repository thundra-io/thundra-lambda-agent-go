package thundraaws

import (
	opentracing "github.com/opentracing/opentracing-go"

	"github.com/aws/aws-sdk-go/aws/request"
	"github.com/aws/aws-sdk-go/aws/session"
)

// Wrap wraps the given session object and adds necessary
// handlers to create a span for the AWS call
func Wrap(s *session.Session) *session.Session {
	if s != nil {
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
	i := integrations[r.ClientInfo.ServiceID]
	span, ctxWithSpan := opentracing.StartSpanFromContext(r.Context(), i.getOperationName(r))
	r.SetContext(ctxWithSpan)
	i.beforeCall(r, span)
}

func completeHandler(r *request.Request) {
	i := integrations[r.ClientInfo.ServiceID]
	span := opentracing.SpanFromContext(r.Context())
	if span == nil {
		return
	}
	i.afterCall(r, span)
	span.Finish()

}

// TODO: Add other handlers