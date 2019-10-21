package tracer

import (
	"encoding/json"
	"os"
	"strings"

	"github.com/thundra-io/thundra-lambda-agent-go/constants"
)

var spanListeners = make([]ThundraSpanListener, 0)

var SpanListenerConstructorMap = make(map[string]func(map[string]interface{}) ThundraSpanListener, 0)

func GetSpanListeners() []ThundraSpanListener {
	return spanListeners
}

func RegisterSpanListener(listener ThundraSpanListener) {
	spanListeners = append(spanListeners, listener)
}

func ClearSpanListeners() {
	spanListeners = make([]ThundraSpanListener, 0)
}

func ParseSpanListeners() {
	ClearSpanListeners()

	for _, env := range os.Environ() {
		if strings.HasPrefix(env, constants.ThundraLambdaSpanListener) {
			config := make(map[string]interface{})
			splits := strings.SplitN(env, "=", 2)

			if err := json.Unmarshal([]byte(splits[1]), &config); err != nil {
				// TODO: Handle json unmarshal error
				continue
			}

			listenerName, ok := config["type"].(string)
			if !ok {
				// TODO: Handle listener type is not string
				continue
			}

			listenerConfig, ok := config["config"].(map[string]interface{})
			if !ok {
				// TODO: Handle config type is not correct
			}

			listenerConstructor, ok := SpanListenerConstructorMap[listenerName]
			if !ok {
				// TODO: Handle listener type does not exist
				continue
			}

			listener := listenerConstructor(listenerConfig)

			if listener != nil {
				RegisterSpanListener(listener)
			}
		}
	}
}

func init() {
	SpanListenerConstructorMap["ErrorInjectorSpanListener"] = NewErrorInjectorSpanListener
	SpanListenerConstructorMap["LatencyInjectorSpanListener"] = NewLatencyInjectorSpanListener
	SpanListenerConstructorMap["FilteringSpanListener"] = NewFilteringSpanListener

	ParseSpanListeners()
}
