package tracer

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

type customError struct {
}

func (e *customError) Error() string {
	return "Custom Error"
}

func TestFrequency(t *testing.T) {
	esl := ErrorInjectorSpanListener{InjectCountFreq: 3}

	for i := 0; i < 10; i++ {
		esl.OnSpanFinished(nil)
	}
	assert.Equal(t, int64(0), esl.counter)

	var errCount int64
	for i := 0; i < 9; i++ {
		func() {
			defer func() {
				if recover() != nil {
					errCount++
				}
			}()
			esl.OnSpanStarted(nil)
		}()

	}
	assert.Equal(t, int64(3), errCount)
	assert.Equal(t, int64(9), esl.counter)
}

func TestError(t *testing.T) {
	errorMessage := "Your name is not good for this mission!"
	esl := ErrorInjectorSpanListener{ErrorMessage: errorMessage, ErrorType: &customError{}}

	var errorPanicked error
	func() {
		defer func() {
			errorPanicked = recover().(error)
		}()
		esl.OnSpanStarted(&spanImpl{})
	}()

	assert.Equal(t, esl.ErrorType, errorPanicked)
	assert.Equal(t, esl.ErrorMessage, errorMessage)

}

func TestNewListenerFromConfig(t *testing.T) {
	config := map[string]interface{}{
		"errorMessage":    "You have a very funny name!",
		"injectOnFinish":  true,
		"injectCountFreq": float64(7),
		"addInfoTags":     false,
		"foo":             "bar",
	}

	esl := NewErrorInjectorSpanListener(config).(*ErrorInjectorSpanListener)

	assert.Equal(t, "You have a very funny name!", esl.ErrorMessage)
	assert.Equal(t, int64(7), esl.InjectCountFreq)
	assert.Equal(t, true, esl.InjectOnFinish)
	assert.Equal(t, false, esl.AddInfoTags)
}

func TestNewListenerFromConfigWithTypeErrors(t *testing.T) {
	config := map[string]interface{}{
		"injectOnFinish":  37,
		"injectCountFreq": "message",
	}

	esl := NewErrorInjectorSpanListener(config).(*ErrorInjectorSpanListener)

	assert.Equal(t, false, esl.InjectOnFinish)
	assert.Equal(t, int64(1), esl.InjectCountFreq)
	assert.Equal(t, true, esl.AddInfoTags)
}
