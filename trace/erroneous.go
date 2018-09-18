package trace

type errorInfo struct {
	Message string
	Kind string
}

type panicInfo struct {
	Message string
	Stack   string
	Kind    string
}
