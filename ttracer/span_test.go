package ttracer

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSetOperationName(t *testing.T) {
	tracer, r := newTracerAndRecorder()

	s := tracer.StartSpan("foo")
	s.SetOperationName("bar")
	s.Finish()

	rs := r.GetSpans()[0]

	assert.True(t, rs.OperationName == "bar")
}

func TestSetTag(t *testing.T) {
	tracer, r := newTracerAndRecorder()

	s := tracer.StartSpan("spiderman")
	s.SetTag("peter", "parker")
	s.SetTag("mary", "jane")
	s.Finish()

	rs := r.GetSpans()[0]

	tags := rs.Tags

	assert.True(t, rs.OperationName == "spiderman")
	assert.True(t, tags["peter"] == "parker")
	assert.True(t, tags["mary"] == "jane")
}

func TestSetParent(t *testing.T) {
	tracer, _ := newTracerAndRecorder()

	ps := tracer.StartSpan("parentSpan")
	cs := tracer.StartSpan("childSpan")

	psi, ok := ps.(*spanImpl)
	assert.True(t, ok)

	csi, ok := cs.(*spanImpl)
	assert.True(t, ok)

	csi.setParent(psi.raw.Context)

	assert.True(t, csi.raw.ParentSpanID == psi.raw.Context.SpanID)
}

func TestLog(t *testing.T) {
	tracer, r := newTracerAndRecorder()

	s := tracer.StartSpan("foo")
	s.LogKV("key1", "value1", "key2", "value2")
	s.Finish()

	rs := r.GetSpans()[0]
	fmt.Println(rs.Logs)

	t.Fail()
}
