package tracer

import (
	"log"
	"strings"
)

type FilteringSpanListener struct {
	Listener ThundraSpanListener
	Filterer SpanFilterer
}

func (f *FilteringSpanListener) OnSpanStarted(span *spanImpl) {
	if f.Listener == nil {
		return
	}

	if f.Filterer == nil || f.Filterer.Accept(span) {
		f.Listener.OnSpanStarted(span)
	}
}

func (f *FilteringSpanListener) OnSpanFinished(span *spanImpl) {
	if f.Listener == nil {
		return
	}

	if f.Filterer == nil || f.Filterer.Accept(span) {
		f.Listener.OnSpanFinished(span)
	}
}

func (f *FilteringSpanListener) PanicOnError() bool {
	return true
}

// NewFilteringSpanListener creates and returns a new FilteringSpanListener from config
func NewFilteringSpanListener(config map[string]string) ThundraSpanListener {
	var listenerStr string

	if config["listener"] != "" {
		listenerStr = config["listener"]
	}

	configPrefix := "config."
	filterPrefix := "filter"
	listenerConfig := make(map[string]string)
	filterConfig := make(map[string]map[string]string)

	for key, val := range config {
		if strings.HasPrefix(key, configPrefix) {

			listenerConfig[key[len(configPrefix):]] = val

		} else if strings.HasPrefix(key, filterPrefix) {

			firstDot := strings.Index(key, ".")
			if firstDot == -1 {
				continue
			}
			filterID := key[:firstDot]
			filterArg := key[firstDot+1:]
			if filterConfig[filterID] == nil {
				filterConfig[filterID] = make(map[string]string, 0)
			}
			filterConfig[filterID][filterArg] = val
		}
	}

	var filters []SpanFilter
	var filterer SpanFilterer
	var listener ThundraSpanListener

	for _, val := range filterConfig {
		filters = append(filters, NewThundraSpanFilter(val))
	}

	if len(filters) > 0 {
		filterer = &ThundraSpanFilterer{filters}
	}

	if SpanListenerConstructorMap[listenerStr] == nil {
		log.Println("No listener found with name:", listenerStr)
	} else {
		listenerConstructor := SpanListenerConstructorMap[listenerStr]
		listener = listenerConstructor(listenerConfig)
	}

	return &FilteringSpanListener{Listener: listener, Filterer: filterer}
}
