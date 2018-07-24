package thundra

import (
	"fmt"
	"os"
	"strconv"
)

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
