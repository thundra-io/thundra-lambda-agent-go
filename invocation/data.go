package invocation

// invocation is the simplest form of data collected from lambda functions. It is collected for any case.
type invocation struct {
	//Base fields
	Id                        string                 `json:"id"`
	Type                      string                 `json:"type"`
	AgentVersion              string                 `json:"agentVersion"`
	DataModelVersion          string                 `json:"dataModelVersion"`
	ApplicationId             string                 `json:"applicationId"`
	ApplicationDomainName     string                 `json:"applicationDomainName"`
	ApplicationClassName      string                 `json:"applicationClassName"`
	ApplicationName           string                 `json:"applicationName"`
	ApplicationVersion        string                 `json:"applicationVersion"`
	ApplicationStage          string                 `json:"applicationStage"`
	ApplicationRuntime        string                 `json:"applicationRuntime"`
	ApplicationRuntimeVersion string                 `json:"applicationRuntimeVersion"`
	ApplicationTags           map[string]interface{} `json:"applicationTags"`

	TraceId          string                 `json:"traceId"`
	TransactionId    string                 `json:"transactionId"`
	SpanId           string                 `json:"spanId"`
	FunctionPlatform string                 `json:"functionPlatform"`
	FunctionName     string                 `json:"functionName"`
	FunctionRegion   string                 `json:"functionRegion"`
	StartTimestamp   int64                  `json:"startTimestamp"`  // Invocation start time in UNIX Epoch milliseconds
	FinishTimestamp  int64                  `json:"finishTimestamp"` // Invocation end time in UNIX Epoch milliseconds
	Duration         int64                  `json:"duration"`        // Invocation time in milliseconds
	Erroneous        bool                   `json:"erroneous"`       // Shows if the invocation failed with an error
	ErrorType        string                 `json:"errorType"`       // Type of the thrown error
	ErrorMessage     string                 `json:"errorMessage"`    // Message of the thrown error
	ErrorCode        string                 `json:"errorCode"`       // Numeric code of the error, such as 404 for HttpError
	ColdStart        bool                   `json:"coldStart"`       // Shows if the invocation is cold started
	Timeout          bool                   `json:"timeout"`         // Shows if the invocation is timed out
	Tags             map[string]interface{} `json:"Tags"`
}
