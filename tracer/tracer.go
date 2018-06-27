package tracer

import (
	ot "github.com/opentracing/opentracing-go"
	"fmt"
	"time"
	"github.com/pkg/errors"
)

type tracer interface {
	ot.Tracer
}

type tracerImpl struct {
}

func (t *tracerImpl) StartSpan(operationName string, opts ...ot.StartSpanOption) ot.Span {
	sso := startSpanOptions{}
	for _, o := range opts {
		o.Apply(&sso.Options)
	}
	// Start time.
	startTime := sso.Options.StartTime
	if startTime.IsZero() {
		startTime = time.Now()
	}

	sp := spanImpl{}

	fmt.Println()
	return nil
}


//
func (tracer *tracerImpl) Inject(sc ot.SpanContext, format interface{}, carrier interface{}) error {
	return errors.New("Inject has not been supported yet")
}

func (tracer *tracerImpl) Extract(format interface{}, carrier interface{}) (ot.SpanContext, error) {
	return nil, errors.New("Extract has not been supported yet")
}
