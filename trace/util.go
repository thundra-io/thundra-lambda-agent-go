package trace

import (
	"fmt"
	"os"
	"strconv"
)

func shouldHideRequest() bool {
	e := os.Getenv(thundraLambdaRequestDisable)
	env, err := strconv.ParseBool(e)
	if err != nil {
		if e != "" {
			fmt.Println(err, thundraLambdaRequestDisable+"is not a bool value. Requests aren't hidden.")
		}
		return false
	}
	return env
}

func shouldHideResponse() bool {
	e := os.Getenv(thundraLambdaResponseDisable)
	env, err := strconv.ParseBool(e)
	if err != nil {
		if e != "" {
			fmt.Println(err, thundraLambdaResponseDisable+" is not a bool value. Responses aren't hidden.")
		}
		return false
	}
	return env
}
