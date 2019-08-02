package config

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetDefaultTimeoutMargin(t *testing.T) {
	AwsLambdaRegion = "us-west-2"
	timeoutMargin := getDefaultTimeoutMargin()
	assert.Equal(t, 200, timeoutMargin)

	AwsLambdaRegion = "us-west-1"
	timeoutMargin = getDefaultTimeoutMargin()
	assert.Equal(t, 400, timeoutMargin)

	AwsLambdaRegion = "us-east-1"
	timeoutMargin = getDefaultTimeoutMargin()
	assert.Equal(t, 600, timeoutMargin)

	AwsLambdaRegion = "eu-west-2"
	timeoutMargin = getDefaultTimeoutMargin()
	assert.Equal(t, 1000, timeoutMargin)

	AwsLambdaRegion = "eu-west-2"
	AwsLambdaFunctionMemorySize = 128
	timeoutMargin = getDefaultTimeoutMargin()
	assert.Equal(t, 3000, timeoutMargin)

	AwsLambdaRegion = "eu-west-2"
	AwsLambdaFunctionMemorySize = 256
	timeoutMargin = getDefaultTimeoutMargin()
	assert.Equal(t, 1500, timeoutMargin)
}
