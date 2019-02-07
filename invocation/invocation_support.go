package invocation

import "fmt"

var invocationTags = make(map[string]interface{})

// SetInvocationTag sets the given tag for invocation
func SetInvocationTag(key string, value interface{}) {
	switch value.(type) {
	case string, int, float64, bool:
		invocationTags[key] = value
	default:
		invocationTags[key] = fmt.Sprint(value)
	}
}

// GetInvocationTags returns invocation tags
func GetInvocationTags() map[string]interface{} {
	return invocationTags
}

// ClearInvocationTags clears the invocation tags
func ClearInvocationTags() {
	invocationTags = make(map[string]interface{})
}
