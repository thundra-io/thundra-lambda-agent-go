package trace

import (
	"context"
	"encoding/json"
	"fmt"
	"testing"
	"time"

	"github.com/opentracing/opentracing-go"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/thundra-io/thundra-lambda-agent-go/plugin"
	"github.com/thundra-io/thundra-lambda-agent-go/test"
	"github.com/thundra-io/thundra-lambda-agent-go/thundra"
)

const (
	mainDuration = 100
	f1Duration   = 50
	f1Name       = "f1"
	f2Duration   = 30
	f2Name       = "f2"
)
// EXAMPLE HANDLERS
func handler1(s string) (string, error) {
	span := opentracing.GlobalTracer().StartSpan("test-operation")
	defer span.Finish()
	span.SetTag("tagKey", "tagValue")
	time.Sleep(time.Millisecond * mainDuration)
	return fmt.Sprintf("Happy monitoring with %s!", s), nil
}

func handler2(s string) (string, error) {
	span := opentracing.GlobalTracer().StartSpan("test-operation")
	defer span.Finish()
	span.SetTag("tagKey", "tagValue")

	ctx := opentracing.ContextWithSpan(context.Background(), span)

	f := func(ctx context.Context, operationName string, duration time.Duration) {
		span, ctx := opentracing.StartSpanFromContext(ctx, operationName)
		defer span.Finish()
		time.Sleep(time.Millisecond * duration)
	}
	f(ctx, "f1", f1Duration)
	f(ctx, "f2", f2Duration)

	time.Sleep(time.Millisecond * mainDuration)
	return fmt.Sprintf("Happy monitoring with %s!", s), nil
}

func TestSpanTransformation(t *testing.T) {

	testCases := []struct {
		name     string
		input    string
		expected expected
		handler  interface{}
	}{
		{
			name:     "Span test with root span only",
			input:    `"Thundra"`,
			expected: expected{"Thundra works!", nil},
			handler:  handler1,
		},
		{
			name:     "Span test with multiple children",
			input:    `"Thundra"`,
			expected: expected{"Thundra works!", nil},
			handler:  handler2,
		},
	}

	for i, testCase := range testCases {
		t.Run(fmt.Sprintf("testCase[%d] %s", i, testCase.name), func(t *testing.T) {
			r := new(test.MockReporter)
			r.On("Report", testApiKey).Return()
			r.On("Clear").Return()
			r.On("Collect", mock.Anything).Return()

			tr := New()
			th := thundra.NewBuilder().AddPlugin(tr).SetReporter(r).SetAPIKey(testApiKey).Build()
			lambdaHandler := thundra.Wrap(testCase.handler, th)
			h := lambdaHandler.(func(context.Context, json.RawMessage) (interface{}, error))
			f := lambdaFunction(h)
			f(context.TODO(), []byte(testCase.input))

			//Monitor Data
			msg, ok := r.MessageQueue[1].(plugin.Message)
			if !ok {
				fmt.Println("Collector message can't be casted to pluginMessage")
			}

			//Trace Data
			td, ok := msg.Data.(traceData)
			if !ok {
				fmt.Println("Can not convert to trace data")
			}

			//Trace Audit Info
			ai := td.AuditInfo
			aiChildren, ok := ai[auditInfoChildren].([]map[string]interface{})
			if !ok {
				fmt.Println("Can not convert auditInfoChildren to []map[string]interface{}")
			}
			root := aiChildren[0]
			assert.Equal(t, "test-operation", root[auditInfoContextName])

			durationMain := root[auditInfoCloseTimestamp].(int64) - root[auditInfoOpenTimestamp].(int64)
			assert.True(t, durationMain >= mainDuration)
			props, ok := root[auditInfoProps].(opentracing.Tags)
			if !ok {
				fmt.Println("auditInfoChildren to opentracing.Tags")
			}
			assert.Equal(t, "tagValue", props["tagKey"])

			if i == 1 {
				secondLevelChildren, _ := root[auditInfoChildren].([]map[string]interface{})
				if !ok {
					fmt.Println("Can not convert to secondLevelChildren")
				}
				assert.Equal(t, f1Name, secondLevelChildren[0][auditInfoContextName])
				assert.Equal(t, f2Name, secondLevelChildren[1][auditInfoContextName])

				// Calculate each functions' duration by endtime - starttime
				duration1 := secondLevelChildren[0][auditInfoCloseTimestamp].(int64) - secondLevelChildren[0][auditInfoOpenTimestamp].(int64)
				duration2 := secondLevelChildren[1][auditInfoCloseTimestamp].(int64) - secondLevelChildren[1][auditInfoOpenTimestamp].(int64)
				assert.True(t, duration1 >= f1Duration)
				assert.True(t, duration2 >= f2Duration)
				assert.True(t, durationMain >= mainDuration+f1Duration+f2Duration)
			}
		})
	}

}