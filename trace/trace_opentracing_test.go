package trace

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"testing"
	"time"

	opentracing "github.com/opentracing/opentracing-go"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/thundra-io/thundra-lambda-agent-go/v2/agent"
	"github.com/thundra-io/thundra-lambda-agent-go/v2/test"
)

const (
	mainDuration = 100
	f1Duration   = 50
	f1Name       = "f1"
	f2Duration   = 30
	f2Name       = "f2"
)

// EXAMPLE HANDLERS
func handler1(ctx context.Context, s string) (string, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "test-operation")
	defer span.Finish()

	span.SetTag("tagKey", "tagValue")
	time.Sleep(time.Millisecond * mainDuration)
	return fmt.Sprintf("Happy monitoring with %s!", s), nil
}

func handler2(ctx context.Context, s string) (string, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "test-operation")
	defer span.Finish()

	span.SetTag("tagKey", "tagValue")

	f := func(ctx context.Context, operationName string, duration time.Duration) {
		data, ctx := opentracing.StartSpanFromContext(ctx, operationName)
		defer data.Finish()
		time.Sleep(time.Millisecond * duration)
	}
	f(ctx, "f1", f1Duration)
	f(ctx, "f2", f2Duration)

	time.Sleep(time.Millisecond * mainDuration)
	return fmt.Sprintf("Happy monitoring with %s!", s), nil
}

func TestSpanTransformation(t *testing.T) {
	// t.Skip("skipping TestSpanTransformation")
	testCases := []struct {
		name     string
		input    string
		expected expected
		handler  interface{}
	}{
		{
			name:     "Span test with root data only",
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
			r := test.NewMockReporter()
			r.On("Report", testAPIKey).Return()
			r.On("Clear").Return()
			r.On("Collect", mock.Anything).Return()

			tr := New()
			a := agent.New().AddPlugin(tr).SetReporter(r)
			lambdaHandler := a.Wrap(testCase.handler)
			h := lambdaHandler.(func(context.Context, json.RawMessage) (interface{}, error))
			f := lambdaFunction(h)
			f(context.TODO(), []byte(testCase.input))

			//Monitor Data
			msg := r.MessageQueue[0]

			// Trace Data
			_, ok := msg.Data.(traceDataModel)
			if !ok {
				log.Println("Can not convert to trace data")
			}

			msg = r.MessageQueue[2]
			// Root Span Data
			rsd, ok := msg.Data.(spanDataModel)
			if !ok {
				log.Println("Can not convert to span data")
			}
			assert.Equal(t, "test-operation", rsd.OperationName)

			durationMain := rsd.Duration
			assert.True(t, durationMain >= mainDuration)

			tags := rsd.Tags
			assert.Equal(t, "tagValue", tags["tagKey"])

			if i == 1 {
				f1Msg := r.MessageQueue[3]
				f2Msg := r.MessageQueue[4]
				// Child span data
				f1Span, ok := f1Msg.Data.(spanDataModel)
				if !ok {
					log.Println("Can not convert f1 span data")
				}

				f2Span, ok := f2Msg.Data.(spanDataModel)
				if !ok {
					log.Println("Can not convert f2 span data")
				}

				assert.Equal(t, "f1", f1Span.OperationName)
				assert.Equal(t, "f2", f2Span.OperationName)
				assert.True(t, f1Span.Duration >= f1Duration)
				assert.True(t, f2Span.Duration >= f2Duration)
			}

		})
	}

}
