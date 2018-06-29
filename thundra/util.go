package thundra

import (
	"fmt"
	"os"
	"strconv"
)

func init() {
	debugEnabled = isThundraDebugEnabled()
}

var debugEnabled bool

func isThundraDisabled() bool {
	env := os.Getenv(thundraLambdaDisable)
	disabled, err := strconv.ParseBool(env)
	if err != nil {
		if env != "" {
			fmt.Println(err, " thundra_lambda_disable is not a bool value. Thundra is enabled by default.")
		}
		return false
	}
	return disabled
}

func isThundraDebugEnabled() bool {
	b, err := strconv.ParseBool(os.Getenv(thundraLambdaDebugEnable))
	if err != nil {
		return false
	}
	return b
}
