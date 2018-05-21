package trace

type errorInfo struct {
	ErrMessage string `json:"errorMessage"`
	ErrType    string `json:"errorType"`
}

type panicInfo struct {
	ErrMessage string `json:"errorMessage"`
	StackTrace string `json:"error"`
	ErrType    string `json:"errorType"`
}
