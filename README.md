# Lambda Go Agent [![CircleCI](https://circleci.com/gh/thundra-io/thundra-lambda-agent-go.svg?style=svg)](https://circleci.com/gh/thundra-io/thundra-lambda-agent-go/) [![Go Report Card](https://goreportcard.com/badge/github.com/thundra-io/thundra-lambda-agent-go)](https://goreportcard.com/report/github.com/thundra-io/thundra-lambda-agent-go)

Trace your AWS lambda functions with async monitoring by [Thundra](https://www.thundra.io/)!

Check out [example projects](https://github.com/thundra-io/thundra-examples-lambda-go) for a quick start and [Thundra docs](https://docs.thundra.io/docs) for more information.

### Usage

In order to trace your lambda usages with Thundra all you need to do is wrap your function.

```
package main

import (
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/thundra-io/thundra-lambda-agent-go/thundra"
	"github.com/thundra-io/thundra-lambda-agent-go/trace"
	"github.com/thundra-io/thundra-lambda-agent-go/metric"
)

//Your lambda handler
func hello() (string, error) {
	return "Hello ƛ!", nil
}

func main() {
	// Instantiate Thundra Agent with Trace & Metric Support
	tr := trace.New()
	m := metric.New()
	t := thundra.NewBuilder().
	            AddPlugin(tr).
	            AddPlugin(m).
	            SetAPIKey(/*TODO login https://console.thundra.io to get your APIKey*/).
	            Build()
	
	// Wrap your lambda function with Thundra
	lambda.Start(thundra.Wrap(hello, t))
}
```
Later just build and deploy your executable to AWS as regular. Test your function on lambda console and visit [Thundra](https://www.thundra.io/) to observe your function metrics.

#### Environment variables

| Name                                     | Type   | Default Value |
|:-----------------------------------------|:------:|:-------------:|
| thundra_apiKey                           | string |       -       |
| thundra_applicationProfile               | string |    default    |
| thundra_disable                          |  bool  |     false     |
| thundra_lambda_trace_request_disable     |  bool  |     false     |
| thundra_lambda_trace_response_disable    |  bool  |     false     |
| thundra_lambda_publish_cloudwatch_enable |  bool  |     false     |
| thundra_lambda_warmup_warmupAware        |  bool  |     false     |
| thundra_lambda_publish_rest_baseUrl      | string |  https<nolink>://collector.thundra.io/api  |
| thundra_log_logLevel                     | string |       -       |
| thundra_lambda_debug_enable              | string |     false     |
| thundra_lambda_timeout_margin            | int    |     180       |

### Async Monitoring

Check out our [docs](https://docs.thundra.io/docs/how-to-setup-async-monitoring) to see how to configure Thundra and async monitoring to visualize your functions in [Thundra](https://www.thundra.io/).

## Warmup Support
You can cut down cold starts easily by deploying our lambda function [`thundra-lambda-warmup`](https://github.com/thundra-io/thundra-lambda-warmup).

Our agent handles warmup requests automatically so you don't need to make any code changes.

You just need to deploy `thundra-lambda-warmup` once, then you can enable warming up for your lambda by 
* setting its environment variable `thundra_lambda_warmup_warmupAware` **true** OR
* adding its name to `thundra-lambda-warmup`'s environment variable `thundra_lambda_warmup_function`.

Check out [this part](https://thundra.readme.io/docs/how-to-warmup) in our docs for more information.