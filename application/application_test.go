package application

import (
	"fmt"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/thundra-io/thundra-lambda-agent-go/constants"
)

func TestParseApplicationTags(t *testing.T) {
	cases := []struct {
		name        string
		key         string
		val         string
		expectedKey string
		expectedVal interface{}
	}{
		{
			name:        "set int application tag",
			key:         constants.ApplicationTagPrefixProp + "intKey",
			val:         "37",
			expectedKey: "intKey",
			expectedVal: int64(37),
		},
		{
			name:        "set bool application tag",
			key:         constants.ApplicationTagPrefixProp + "boolKey",
			val:         "true",
			expectedKey: "boolKey",
			expectedVal: true,
		},
		{
			name:        "set bool application tag",
			key:         constants.ApplicationTagPrefixProp + "boolKey",
			val:         "t",
			expectedKey: "boolKey",
			expectedVal: true,
		},
		{
			name:        "set bool application tag",
			key:         constants.ApplicationTagPrefixProp + "boolKey",
			val:         "1",
			expectedKey: "boolKey",
			expectedVal: true,
		},
		{
			name:        "set bool application tag",
			key:         constants.ApplicationTagPrefixProp + "boolKey",
			val:         "false",
			expectedKey: "boolKey",
			expectedVal: false,
		},
		{
			name:        "set bool application tag",
			key:         constants.ApplicationTagPrefixProp + "boolKey",
			val:         "f",
			expectedKey: "boolKey",
			expectedVal: false,
		},
		{
			name:        "set bool application tag",
			key:         constants.ApplicationTagPrefixProp + "boolKey",
			val:         "0",
			expectedKey: "boolKey",
			expectedVal: false,
		},
		{
			name:        "set float application tag",
			key:         constants.ApplicationTagPrefixProp + "floatKey",
			val:         "1.5",
			expectedKey: "floatKey",
			expectedVal: 1.5,
		},
		{
			name:        "set float application tag",
			key:         constants.ApplicationTagPrefixProp + "stringKey",
			val:         "foobar",
			expectedKey: "stringKey",
			expectedVal: "foobar",
		},
	}

	for i, testCase := range cases {
		t.Run(fmt.Sprintf("testCase[%d] %s", i, testCase.name), func(t *testing.T) {
			os.Setenv(testCase.key, testCase.val)
			parseApplicationTags()
			assert.Equal(t, testCase.expectedVal, ApplicationTags[testCase.expectedKey])
			os.Unsetenv(testCase.key)
		})
	}
}

func TestApplicationDomainNameFromEnv(t *testing.T) {
	os.Setenv(constants.ApplicationDomainProp, "fooDomain")
	domainName := getApplicationDomainName()
	assert.Equal(t, domainName, "fooDomain")
	os.Unsetenv(constants.ApplicationDomainProp)
}

func TestApplicationDomainNameDefault(t *testing.T) {
	domainName := getApplicationDomainName()
	assert.Equal(t, domainName, constants.AwsLambdaApplicationDomain)
}

func TestApplicationClassNameFromEnv(t *testing.T) {
	os.Setenv(constants.ApplicationClassProp, "fooClass")
	className := getApplicationClassName()
	assert.Equal(t, className, "fooClass")
	os.Unsetenv(constants.ApplicationClassProp)
}

func TestApplicationClassNameDefault(t *testing.T) {
	className := getApplicationClassName()
	assert.Equal(t, className, constants.AwsLambdaApplicationClass)
}

func TestApplicationStageFromEnv(t *testing.T) {
	os.Setenv(constants.ApplicationStageProp, "fooStage")
	stage := getApplicationStage()
	assert.Equal(t, stage, "fooStage")
	os.Unsetenv(constants.ApplicationStageProp)
}

func TestApplicationStageDefault(t *testing.T) {
	stage := getApplicationStage()
	assert.Equal(t, stage, os.Getenv(constants.ThundraApplicationStage))
}
