package trace

import (
	"fmt"
	"os"
	"strconv"
)

func shouldHideRequest() bool {
	e := os.Getenv(thundraLambdaTraceRequestDisable)
	env, err := strconv.ParseBool(e)
	if err != nil {
		if e != "" {
			fmt.Println(err, thundraLambdaTraceRequestDisable+"is not a bool value. Requests aren't hidden.")
		}
		return false
	}
	return env
}

func shouldHideResponse() bool {
	e := os.Getenv(thundraLambdaTraceResponseDisable)
	env, err := strconv.ParseBool(e)
	if err != nil {
		if e != "" {
			fmt.Println(err, thundraLambdaTraceResponseDisable+" is not a bool value. Responses aren't hidden.")
		}
		return false
	}
	return env
}
