package tracer

import (
	"log"
	"os"
	"regexp"
	"strings"

	"github.com/thundra-io/thundra-lambda-agent-go/constants"
)

var spanListeners = make([]ThundraSpanListener, 0)

var SpanListenerConstructorMap = make(map[string]func(map[string]string) ThundraSpanListener, 0)

var p, _ = regexp.Compile(`(\w+)\[(.*)\]`)
var p2, _ = regexp.Compile(`[\w.]+=[\w.\!\?\-\"\'\:\(\) ]*`)

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
			splits := strings.SplitN(env, "=", 2)
			listenerStr, configStr := getListenerAndConfig(splits[1])
			if listenerStr != "" {
				if SpanListenerConstructorMap[listenerStr] == nil {
					log.Println("No listener found with name:", listenerStr)
				} else {
					listenerContructor := SpanListenerConstructorMap[listenerStr]
					config := parseConfig(configStr)
					listener := listenerContructor(config)
					RegisterSpanListener(listener)
				}
			}
		}
	}
}

func parseConfig(configStr string) map[string]string {
	config := make(map[string]string, 0)

	if configStr != "" {
		res := p2.FindAllString(configStr, -1)
		for i := range res {
			pair := strings.Split(res[i], "=")
			if len(pair) == 2 {
				config[pair[0]] = pair[1]
			}
		}
	}
	return config
}

func getListenerAndConfig(envValue string) (string, string) {
	res := p.FindAllStringSubmatch(envValue, -1)
	var listener, config = "", ""
	for i := range res {
		listener = res[i][1]
		config = res[i][2]
	}
	return listener, config
}

func init() {
	SpanListenerConstructorMap["ErrorInjectorSpanListener"] = NewErrorInjectorSpanListener
	SpanListenerConstructorMap["LatencyInjectorSpanListener"] = NewLatencyInjectorSpanListener
	SpanListenerConstructorMap["FilteringSpanListener"] = NewFilteringSpanListener

	ParseSpanListeners()
}
