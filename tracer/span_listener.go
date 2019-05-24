package tracer

type ThundraSpanListener interface {
	OnSpanStarted(*spanImpl)
	OnSpanFinished(*spanImpl)
	PanicOnError() bool
}
