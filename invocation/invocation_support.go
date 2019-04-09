package invocation

var invocationTags = make(map[string]interface{})

// SetTag sets the given tag for invocation
func SetTag(key string, value interface{}) {
	invocationTags[key] = value
}

// GetTag returns invocation tag for key
func GetTag(key string) interface{} {
	return invocationTags[key]
}

// GetTags returns invocation tags
func GetTags() map[string]interface{} {
	return invocationTags
}

// ClearTags clears the invocation tags
func ClearTags() {
	invocationTags = make(map[string]interface{})
}
