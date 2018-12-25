package ttracer

type SpanEvent int

const (
	StartSpanEvent  SpanEvent = 0
	FinishSpanEvent SpanEvent = 1
)
