package apexgatewayv2

import (
	"context"
	"github.com/apex/gateway/v2"
	"github.com/thundra-io/thundra-lambda-agent-go/v2/thundra"
	"net/http"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

func ListenAndServe(h http.Handler) error {
	if h == nil {
		h = http.DefaultServeMux
	}

	lambda.Start(thundra.Wrap(wrapper(h)))

	return nil
}

func wrapper(h http.Handler) func(ctx context.Context, e events.APIGatewayV2HTTPRequest) (events.APIGatewayV2HTTPResponse, error) {
	return func(ctx context.Context, e events.APIGatewayV2HTTPRequest) (events.APIGatewayV2HTTPResponse, error) {
		r, err := gateway.NewRequest(ctx, e)
		if err != nil {
			return events.APIGatewayV2HTTPResponse{}, err
		}

		w := gateway.NewResponse()
		h.ServeHTTP(w, r)

		return w.End(), nil
	}
}
