package tracer

import (
	"reflect"
	"strconv"
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
		"error_type":        reflect.TypeOf(err),
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
func NewErrorInjectorSpanListener(config map[string]string) ThundraSpanListener {

	spanListener := &ErrorInjectorSpanListener{ErrorMessage: defaultErrorMessage, AddInfoTags: true, InjectCountFreq: 1}

	if config["errorMessage"] != "" {
		spanListener.ErrorMessage = config["errorMessage"]
	}
	if config["injectOnFinish"] != "" {
		injectOnFinish, err := strconv.ParseBool(config["injectOnFinish"])
		if err == nil {
			spanListener.InjectOnFinish = injectOnFinish
		}
	}
	if config["injectCountFreq"] != "" {
		injectCountFreq, err := strconv.ParseInt(config["injectCountFreq"], 10, 64)
		if err == nil {
			spanListener.InjectCountFreq = injectCountFreq
		}
	}
	if config["addInfoTags"] != "" {
		addInfoTags, err := strconv.ParseBool(config["addInfoTags"])
		if err == nil {
			spanListener.AddInfoTags = addInfoTags
		}
	}
	spanListener.ErrorType = errors.New(spanListener.ErrorMessage)

	return spanListener
}
