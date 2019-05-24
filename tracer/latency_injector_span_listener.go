package tracer

import (
	"math/rand"
	"strconv"
	"time"

	"github.com/thundra-io/thundra-lambda-agent-go/constants"
)

var defaultDelay int64 = 100

type LatencyInjectorSpanListener struct {
	Delay          int64
	InjectOnFinish bool
	RandomizeDelay bool
	AddInfoTags    bool
}

func (l *LatencyInjectorSpanListener) OnSpanStarted(span *spanImpl) {
	if !l.InjectOnFinish {
		l.injectDelay(span)
	}
}

func (l *LatencyInjectorSpanListener) OnSpanFinished(span *spanImpl) {
	if l.InjectOnFinish {
		l.injectDelay(span)
	}
}

func (l *LatencyInjectorSpanListener) PanicOnError() bool {
	return false
}

func (l *LatencyInjectorSpanListener) injectDelay(span *spanImpl) {
	delay := l.Delay
	if delay <= 0 {
		delay = defaultDelay
	}
	if l.RandomizeDelay {
		delay = rand.Int63n(delay)
	}
	if l.AddInfoTags {
		l.addInfoTags(span, delay)
	}
	time.Sleep(time.Duration(delay) * time.Millisecond)
}

func (l *LatencyInjectorSpanListener) addInfoTags(span *spanImpl, injectedDelay int64) {
	infoTags := map[string]interface{}{
		"type":             "latency_injecter_span_listener",
		"inject_on_finish": l.InjectOnFinish,
		"delay":            l.Delay,
		"injected_delay":   injectedDelay,
	}
	span.SetTag(constants.ThundraLambdaSpanListenerInfoTag, infoTags)
}

// NewLatencyInjectorSpanListener creates and returns a new LatencyInjectorSpanListener from config
func NewLatencyInjectorSpanListener(config map[string]string) ThundraSpanListener {
	spanListener := &LatencyInjectorSpanListener{Delay: defaultDelay, AddInfoTags: true}

	if config["injectOnFinish"] != "" {
		injectOnFinish, err := strconv.ParseBool(config["injectOnFinish"])
		if err == nil {
			spanListener.InjectOnFinish = injectOnFinish
		}

	}
	if config["delay"] != "" {
		delay, err := strconv.ParseInt(config["delay"], 10, 64)
		if err == nil {
			spanListener.Delay = delay
		}
	}
	if config["addInfoTags"] != "" {
		addInfoTags, err := strconv.ParseBool(config["addInfoTags"])
		if err == nil {
			spanListener.AddInfoTags = addInfoTags
		}

	}
	if config["randomizeDelay"] != "" {
		randomizeDelay, err := strconv.ParseBool(config["randomizeDelay"])
		if err == nil {
			spanListener.RandomizeDelay = randomizeDelay
		}
	}

	return spanListener
}
