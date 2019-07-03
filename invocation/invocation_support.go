package invocation

var invocationTags = make(map[string]interface{})
var userInvocationTags = make(map[string]interface{})

// SetTag sets the given tag for invocation
func SetTag(key string, value interface{}) {
	userInvocationTags[key] = value
}

// GetTag returns invocation tag for key
func GetTag(key string) interface{} {
	return userInvocationTags[key]
}

// GetTags returns invocation tags
func GetTags() map[string]interface{} {
	return userInvocationTags
}

// SetAgentTag sets the given tag for invocation
func SetAgentTag(key string, value interface{}) {
	invocationTags[key] = value
}

// GetAgentTag returns invocation tag for key
func GetAgentTag(key string) interface{} {
	return invocationTags[key]
}

// GetAgentTags returns invocation tags
func GetAgentTags() map[string]interface{} {
	return invocationTags
}

// ClearTags clears the invocation tags
func ClearTags() {
	invocationTags = make(map[string]interface{})
	userInvocationTags = make(map[string]interface{})
}
