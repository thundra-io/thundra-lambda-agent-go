package tracer

import (
	"bytes"
	"compress/gzip"
	"encoding/base64"
	"encoding/json"
	"io/ioutil"
	"log"
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

			if len(splits) < 2 {
				continue
			}

			var err error
			configStr := splits[1]

			if !strings.HasPrefix(configStr, "{") {
				configStr, err = decodeConfigStr(configStr)
				if err != nil {
					log.Println("Couldn't parse given span listener configuration:", err)
					continue
				}
			}

			if err := json.Unmarshal([]byte(configStr), &config); err != nil {
				log.Println("Given span listener configuration is not a valid JSON string:", err)
				continue
			}

			listenerName, ok := config["type"].(string)
			if !ok {
				log.Println("Given listener type is not a valid span listener")
				continue
			}

			listenerConfig, ok := config["config"].(map[string]interface{})
			if !ok {
				log.Println("No config given for the span listener")
			}

			listenerConstructor, ok := SpanListenerConstructorMap[listenerName]
			if !ok {
				log.Println("Given listener type is not a valid span listener")
				continue
			}

			listener := listenerConstructor(listenerConfig)

			if listener != nil {
				RegisterSpanListener(listener)
			}
		}
	}
}

func decodeConfigStr(configStr string) (string, error) {
	z, err := base64.StdEncoding.DecodeString(configStr)
	if err != nil {
		return "", err
	}

	r, err := gzip.NewReader(bytes.NewReader(z))
	if err != nil {
		return "", err
	}

	result, err := ioutil.ReadAll(r)
	if err != nil {
		return "", err
	}

	return string(result), nil
}

func init() {
	SpanListenerConstructorMap["ErrorInjectorSpanListener"] = NewErrorInjectorSpanListener
	SpanListenerConstructorMap["LatencyInjectorSpanListener"] = NewLatencyInjectorSpanListener
	SpanListenerConstructorMap["FilteringSpanListener"] = NewFilteringSpanListener
	SpanListenerConstructorMap["TagInjectorSpanListener"] = NewTagInjectorSpanListener
	SpanListenerConstructorMap["SecurityAwareSpanListener"] = NewSecurityAwareSpanListener
	ParseSpanListeners()
}
