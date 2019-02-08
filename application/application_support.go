package application

import (
	"os"
	"strconv"
	"strings"

	"github.com/thundra-io/thundra-lambda-agent-go/constants"
)

var applicationTags = make(map[string]interface{})

func parseApplicationTags() {
	tagPrefix := constants.ApplicationTagPrefixProp
	prefixLen := len(tagPrefix)
	for _, pair := range os.Environ() {
		if strings.HasPrefix(pair, tagPrefix) {
			splits := strings.SplitN(pair[prefixLen:], "=", 2)
			key, val := splits[0], splits[1]
			applicationTags[key] = parseStringToVal(val)
		}
	}
}

func parseStringToVal(s string) interface{} {
	if v, err := strconv.ParseBool(s); err == nil {
		return v
	}
	if v, err := strconv.ParseInt(s, 10, 32); err == nil {
		return v
	}
	if v, err := strconv.ParseFloat(s, 32); err == nil {
		return v
	}
	return strings.Trim(s, "\"")
}
