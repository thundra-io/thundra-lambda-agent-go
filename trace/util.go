package trace

import (
	"fmt"
	"os"
	"strconv"
)

func shouldHideRequest() bool {
	e := os.Getenv(thundraLambdaHideRequest)
	env, err := strconv.ParseBool(e)
	if err != nil {
		if e != "" {
			fmt.Println(err, " thundra_lambda_hide_request is not a bool value. Requests aren't hidden.")
		}
		return false
	}
	return env
}

func shouldHideResponse() bool {
	e := os.Getenv(thundraLambdaHideResponse)
	env, err := strconv.ParseBool(e)
	if err != nil {
		if e != "" {
			fmt.Println(err, " thundra_lambda_hide_response is not a bool value. Responses aren't hidden.")
		}
		return false
	}
	return env
}
