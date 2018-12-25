package ttracer

// Options allows creating a customized Tracer via NewWithOptions.
type Options struct {
	// MaxLogsPerSpan limits the number of Logs in a span (if set to a nonzero
	// value). If a span has more logs than this value, logs are dropped as
	// necessary (and replaced with a log describing how many were dropped).
	//
	// About half of the MaxLogPerSpan logs kept are the oldest logs, and about
	// half are the newest logs.
	//
	// This value is ignored if DropAllLogs is true.
	MaxLogsPerSpan int
	// Recorder receives Spans when they have been started or finished.
	Recorder SpanRecorder
	// DropAllLogs turns log events on all Spans into no-ops.
	DropAllLogs bool
}

// DefaultOptions returns an Options object with all options disabled.
// A Recorder needs to be set manually before using the returned object with a Tracer.
func DefaultOptions() Options {
	return Options{
		MaxLogsPerSpan: 100,
	}
}
