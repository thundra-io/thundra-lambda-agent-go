# Lambda Go Agent

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
)

//Your lambda handler
func hello() (string, error) {
	return "Hello ƛ!", nil
}

func main() {
	// Instantiate Thundra Agent with Trace Support
	t := thundra.NewBuilder().
	            AddPlugin(&trace.Trace{}).
	            SetAPIKey(/*TODO login https://console.thundra.io to get your APIKey*/).
	            Build()
	
	// Wrap your lambda function with Thundra
	lambda.Start(thundra.Wrap(hello, t))
}
```
Later just build and deploy your executable to AWS as regular. Test your function on lambda console and visit [Thundra](https://www.thundra.io/) to observe your function metrics.

### Async Monitoring

Check out our [docs](https://docs.thundra.io/docs/how-to-setup-async-monitoring) to see how to configure Thundra and async monitoring to visualize your functions in [Thundra](https://www.thundra.io/).
