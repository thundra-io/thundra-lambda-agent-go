package plugin

import (
	"github.com/satori/go.uuid"
	"strings"
)

func GenerateNewId() string {
	return uuid.Must(uuid.NewV4()).String()
}

func SplitAppId(logStreamName string) string {
	s := strings.Split(logStreamName, "]")
	if len(s) > 1 {
		return s[1]
	} else {
		return ""
	}
}
