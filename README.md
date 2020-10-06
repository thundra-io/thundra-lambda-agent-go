# Lambda Go Agent [![OpenTracing Badge](https://img.shields.io/badge/OpenTracing-enabled-blue.svg)](http://opentracing.io) [![CircleCI](https://circleci.com/gh/thundra-io/thundra-lambda-agent-go.svg?style=svg)](https://circleci.com/gh/thundra-io/thundra-lambda-agent-go/) [![Go Report Card](https://goreportcard.com/badge/github.com/thundra-io/thundra-lambda-agent-go)](https://goreportcard.com/report/github.com/thundra-io/thundra-lambda-agent-go)

Trace your AWS lambda functions with async monitoring by [Thundra](https://www.thundra.io/)!

Check out [Thundra docs](https://docs.thundra.io/docs) for more information.

### Usage

In order to trace your lambda usages with Thundra all you need to do is wrap your function.
```go
import "github.com/thundra-io/thundra-lambda-agent-go/v2/thundra"	// with go modules enabled (GO111MODULE=on or outside GOPATH) for version >= v2.3.1
import "github.com/thundra-io/thundra-lambda-agent-go/thundra"         // with go modules disabled
```

```go
package main

import (
	"github.com/aws/aws-lambda-go/lambda"
	// thundra import here
)

// Your lambda handler
func handler() (string, error) {
	return "Hello, Thundra!", nil
}

func main() {
	// Wrap your lambda handler with Thundra
	lambda.Start(thundra.Wrap(handler))
}
```
Later just build and deploy your executable to AWS as regular. Test your function on lambda console and visit [Thundra](https://console.thundra.io/) to observe your function metrics.

#### Environment variables

| Name                                                  |   Type     |    Default Value                |
|:------------------------------------------------------|:----------:|:-------------------------------:|
| thundra_applicationProfile                            |   string   |    default                      |
| thundra_agent_lambda_disable                          |   bool     |    false                        |
| thundra_agent_lambda_timeout_margin                   |   number   |    200                          |
| thundra_agent_lambda_report_rest_baseUrl              |   string   |    https://api.thundra.io/v1    |
| thundra_agent_lambda_trace_disable                    |   bool     |    false                        |
| thundra_agent_lambda_metric_disable                   |   bool     |    false                        |
| thundra_agent_lambda_log_disable                      |   bool     |    false                        |
| thundra_log_logLevel                                  |   string   |    TRACE                        |
| thundra_agent_lambda_trace_request_skip               |   bool     |    false                        |
| thundra_agent_lambda_trace_response_skip              |   bool     |    false                        |
| thundra_agent_lambda_report_rest_trustAllCertificates |   bool     |    false                        |
| thundra_agent_lambda_debug_enable                     |   bool     |    false                        |
| thundra_agent_lambda_warmup_warmupAware               |   bool     |    false                        |


### Async Monitoring

Check out our [docs](https://docs.thundra.io/docs/how-to-setup-async-monitoring) to see how to configure Thundra and async monitoring to visualize your functions in [Thundra](https://www.thundra.io/).

## Warmup Support
You can cut down cold starts easily by deploying our lambda function [`thundra-lambda-warmup`](https://github.com/thundra-io/thundra-lambda-warmup).

Our agent handles warmup requests automatically so you don't need to make any code changes.

You just need to deploy `thundra-lambda-warmup` once, then you can enable warming up for your lambda by 
* setting its environment variable `thundra_agent_lambda_warmup_warmupAware` **true** OR
* adding its name to `thundra-lambda-warmup`'s environment variable `thundra_lambda_warmup_function`.

Check out [this part](https://thundra.readme.io/docs/how-to-warmup) in our docs for more information.
