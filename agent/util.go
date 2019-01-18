package agent

import (
	"context"
	"encoding/json"
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

func (handler lambdaFunction) invoke(ctx context.Context, payload []byte) ([]byte, error) {
	response, err := handler(ctx, payload)
	if err != nil {
		return nil, err
	}

	responseBytes, err := json.Marshal(response)
	if err != nil {
		return nil, err
	}

	return responseBytes, nil
}

type expected struct {
	val string
	err error
}
