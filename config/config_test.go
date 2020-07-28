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
	regions := []string{
		"us-west-2", "us-west-1",
		"us-east-2", "us-east-1",
		"ca-central-1", "sa-east-1",
		"eu-central-1", "eu-west-1",
		"eu-west-2", "eu-west-3",
		"eu-north-1", "eu-south-1",
		"ap-south-1", "ap-northeast-1",
		"ap-northeast-2", "ap-southeast-1",
		"ap-southeast-2", "ap-east-1",
		"af-south-1", "me-south-1",
	}

	for _, region := range regions {
		AwsLambdaRegion = region
		collector := getDefaultCollector()
		assert.Equal(t, region+".collector.thundra.io", collector)
	}

	AwsLambdaRegion = ""
	collector := getDefaultCollector()
	assert.Equal(t, "collector.thundra.io", collector)
}
