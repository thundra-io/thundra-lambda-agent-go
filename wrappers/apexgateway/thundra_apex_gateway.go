package apexgateway

import (
	"context"
	"github.com/apex/gateway"
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

func wrapper(h http.Handler) func(ctx context.Context, e events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	return func(ctx context.Context, e events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
		r, err := gateway.NewRequest(ctx, e)
		if err != nil {
			return events.APIGatewayProxyResponse{}, err
		}

		w := gateway.NewResponse()
		h.ServeHTTP(w, r)
		return w.End(), nil
	}
}
