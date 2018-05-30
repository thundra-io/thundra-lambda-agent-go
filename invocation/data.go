package invocation

type invocation struct {
	Id                 string `json:"id"`
	TransactionId      string `json:"transactionId"`
	ApplicationName    string `json:"applicationName"`
	ApplicationId      string `json:"applicationId"`
	ApplicationVersion string `json:"applicationVersion"`
	ApplicationProfile string `json:"applicationProfile"`
	ApplicationType    string `json:"applicationType"`

	Duration       int64  `json:"duration"`       // Invocation time in milliseconds
	StartTimestamp int64  `json:"startTimestamp"` // Invocation start time in UNIX Epoch milliseconds
	EndTimestamp   int64  `json:"endTimestamp"`   // Invocation end time in UNIX Epoch milliseconds
	Erroneous      bool   `json:"erroneous"`      // Shows if the invocation failed with an error
	ErrorType      string `json:"errorType"`      // Type of the thrown error
	ErrorMessage   string `json:"errorMessage"`   // Message of the thrown error
	ColdStart      bool   `json:"coldStart"`      // Shows if the invocation is cold started
	Timeout        bool   `json:"timeout"`        // Shows if the invocation is timed out

	Region     string `json:"region"`     // Name of the AWS region
	MemorySize int    `json:"memorySize"` // Memory Size of the function in MB
}
