package invocation

var invocationTags = make(map[string]interface{})

// SetTag sets the given tag for invocation
func SetTag(key string, value interface{}) {
	invocationTags[key] = value
}

// GetTags returns invocation tags
func GetTags() map[string]interface{} {
	return invocationTags
}

// ClearTags clears the invocation tags
func ClearTags() {
	invocationTags = make(map[string]interface{})
}
