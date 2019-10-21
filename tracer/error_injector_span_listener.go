package tracer

import (
	"log"
	"reflect"
	"sync/atomic"

	"github.com/pkg/errors"
	"github.com/thundra-io/thundra-lambda-agent-go/constants"
	"github.com/thundra-io/thundra-lambda-agent-go/utils"
)

var defaultErrorMessage = "Error injected by Thundra!"

type ErrorInjectorSpanListener struct {
	ErrorMessage    string
	ErrorType       error
	InjectOnFinish  bool
	InjectCountFreq int64
	counter         int64
	AddInfoTags     bool
}

func (e *ErrorInjectorSpanListener) OnSpanStarted(span *spanImpl) {

	if !e.InjectOnFinish && e.ableToRaise() {
		e.injectError(span)
	}
}

func (e *ErrorInjectorSpanListener) OnSpanFinished(span *spanImpl) {
	if e.InjectOnFinish && e.ableToRaise() {
		e.injectError(span)
	}
}

func (e *ErrorInjectorSpanListener) PanicOnError() bool {
	return true
}

func (e *ErrorInjectorSpanListener) ableToRaise() bool {
	counter := atomic.AddInt64(&e.counter, 1)
	countfreq := e.InjectCountFreq
	if e.InjectCountFreq < 1 {
		countfreq = 1
	}
	return (counter % countfreq) == 0
}

func (e *ErrorInjectorSpanListener) addInfoTags(span *spanImpl, err error) {
	infoTags := map[string]interface{}{
		"type":              "error_injecter_span_listener",
		"error_type":        reflect.TypeOf(err).String(),
		"error_message":     err.Error(),
		"inject_on_finish":  e.InjectOnFinish,
		"inject_count_freq": e.InjectCountFreq,
	}
	span.SetTag(constants.ThundraLambdaSpanListenerInfoTag, infoTags)
}

func (e *ErrorInjectorSpanListener) injectError(span *spanImpl) {
	var err error
	var errMessage = defaultErrorMessage

	if e.ErrorMessage != "" {
		errMessage = e.ErrorMessage
	}

	if e.ErrorType != nil {
		err = e.ErrorType
	} else {
		err = errors.New(errMessage)
	}
	utils.SetSpanError(span, err)
	if e.AddInfoTags {
		e.addInfoTags(span, err)
	}
	panic(err)
}

// NewErrorInjectorSpanListener creates and returns a new ErrorInjectorSpanListener from config
func NewErrorInjectorSpanListener(config map[string]interface{}) ThundraSpanListener {
	spanListener := &ErrorInjectorSpanListener{ErrorMessage: defaultErrorMessage, AddInfoTags: true, InjectCountFreq: 1}

	if errorMessage, ok := config["errorMessage"].(string); ok {
		spanListener.ErrorMessage = errorMessage
	}
	if injectOnFinish, ok := config["injectOnFinish"].(bool); ok {
		spanListener.InjectOnFinish = injectOnFinish
	}
	if injectCountFreq, ok := config["injectCountFreq"].(int64); ok {
		spanListener.InjectCountFreq = injectCountFreq
	}
	if addInfoTags, ok := config["addInfoTags"].(bool); ok {
		spanListener.AddInfoTags = addInfoTags
	}
	spanListener.ErrorType = errors.New(spanListener.ErrorMessage)

	return spanListener
}
