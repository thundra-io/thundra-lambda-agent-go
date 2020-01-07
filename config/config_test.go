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

func TestGetNearestCollector(t *testing.T) {
	AwsLambdaRegion = "us-west-1"
	collector := getNearestCollector()
	assert.Equal(t, "api.thundra.io", collector)

	AwsLambdaRegion = "us-east-1"
	collector = getNearestCollector()
	assert.Equal(t, "api-us-east-1.thundra.io", collector)

	AwsLambdaRegion = "eu-west-2"
	collector = getNearestCollector()
	assert.Equal(t, "api-eu-west-2.thundra.io", collector)

	AwsLambdaRegion = "eu-west-1"
	collector = getNearestCollector()
	assert.Equal(t, "api-eu-west-1.thundra.io", collector)

	AwsLambdaRegion = "ap-"
	collector = getNearestCollector()
	assert.Equal(t, "api-ap-northeast-1.thundra.io", collector)
}
