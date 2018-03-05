package thundra

import (
	"testing"
	"fmt"
	"context"
	"errors"
	"github.com/stretchr/testify/assert"
	"encoding/json"
)

// Invoke calls the handler, and serializes the response.
// If the underlying handler returned an error, or an error occurs during serialization, error is returned.
func (handler ThundraLambdaHandler) Invoke(ctx context.Context, payload []byte) ([]byte, error) {
	response, err := handler(ctx, payload)
	if err != nil {
		return nil, err
	}

	responseBytes, err := json.Marshal(response)
	if err != nil {
		return nil, err
	}

	return responseBytes, nil
}

type expected struct {
	val string
	err error
}

func TestWrapper(t *testing.T) {
	hello := func(s string) string {
		return fmt.Sprintf("Happy monitoring with %s!", s)
	}
	hellop := func(s *string) *string {
		v := hello(*s)
		return &v
	}

	testCases := []struct {
		name     string
		input    string
		expected expected
		handler  interface{}
	}{
		{
			input:    `"Thundra"`,
			expected: expected{`"Happy monitoring with Thundra!"`, nil},
			handler: func(name string) (string, error) {
				return hello(name), nil
			},
		},
		{
			input:    `"Thundra"`,
			expected: expected{`"Happy monitoring with Thundra!"`, nil},
			handler: func(name string) (string, error) {
				return hello(name), nil
			},
		},
		{
			input:    `"Thundra"`,
			expected: expected{`"Happy monitoring with Thundra!"`, nil},
			handler: func(ctx context.Context, name string) (string, error) {
				return hello(name), nil
			},
		},
		{
			input:    `"Thundra"`,
			expected: expected{`"Happy monitoring with Thundra!"`, nil},
			handler: func(name *string) (*string, error) {
				return hellop(name), nil
			},
		},
		{
			input:    `"Thundra"`,
			expected: expected{`"Happy monitoring with Thundra!"`, nil},
			handler: func(name *string) (*string, error) {
				return hellop(name), nil
			},
		},
		{
			input:    `"Thundra"`,
			expected: expected{`"Happy monitoring with Thundra!"`, nil},
			handler: func(ctx context.Context, name *string) (*string, error) {
				return hellop(name), nil
			},
		},
		{
			input:    `"Thundra"`,
			expected: expected{"", errors.New("thundra is dead baby, thundra is dead")},
			handler: func() error {
				return errors.New("thundra is dead baby, thundra is dead")
			},
		},
		{
			input:    `"Thundra"`,
			expected: expected{"", errors.New("thundra is dead baby, thundra is dead")},
			handler: func() (interface{}, error) {
				return nil, errors.New("thundra is dead baby, thundra is dead")
			},
		},
		{
			input:    `"Thundra"`,
			expected: expected{"", errors.New("thundra is dead baby, thundra is dead")},
			handler: func(e interface{}) (interface{}, error) {
				return nil, errors.New("thundra is dead baby, thundra is dead")
			},
		},
		{
			input:    `"Thundra"`,
			expected: expected{"", errors.New("thundra is dead baby, thundra is dead")},
			handler: func(ctx context.Context, e interface{}) (interface{}, error) {
				return nil, errors.New("thundra is dead baby, thundra is dead")
			},
		},
		{
			name:     "basic input struct serialization",
			input:    `{"custom":9001}`,
			expected: expected{`9001`, nil},
			handler: func(event struct{ Custom int }) (int, error) {
				return event.Custom, nil
			},
		},
		{
			name:     "basic output struct serialization",
			input:    `9001`,
			expected: expected{`{"Number":9001}`, nil},
			handler: func(event int) (struct{ Number int }, error) {
				return struct{ Number int }{event}, nil
			},
		},
	}
	for i, testCase := range testCases {
		t.Run(fmt.Sprintf("testCase[%d] %s", i, testCase.name), func(t *testing.T) {
			th := GetInstance([]string{})

			lambdaHandler := WrapLambdaHandler(testCase.handler, th)
			response, err := lambdaHandler.Invoke(context.TODO(), []byte(testCase.input))

			if testCase.expected.err != nil {
				assert.Equal(t, testCase.expected.err, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, testCase.expected.val, string(response))
			}
		})
	}
}