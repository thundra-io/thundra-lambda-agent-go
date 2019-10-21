package tracer

import "log"

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
func NewFilteringSpanListener(config map[string]interface{}) ThundraSpanListener {
	listenerDef, ok := config["listener"].(map[string]interface{})
	log.Println(listenerDef)
	if !ok {
		// TODO: Handle listener definition type error
		return nil
	}

	listenerConfig, ok := listenerDef["config"].(map[string]interface{})
	log.Println(listenerConfig)
	if !ok {
		// TODO: Handle listener config type
	}

	all, ok := config["all"].(bool)
	if !ok {
		// TODO: Handle all value is not bool
	}

	filterer := &ThundraSpanFilterer{spanFilters: []SpanFilter{}, all: all}

	if filterConfigs, ok := config["filters"].([]interface{}); ok {
		filterer.spanFilters = crateFiltersFromConfig(filterConfigs)
	}

	listenerName, ok := listenerDef["type"].(string)
	log.Println(listenerName)
	if !ok {
		// TODO: Handle listener type name
	}

	listenerConstructor, ok := SpanListenerConstructorMap[listenerName]
	if !ok {
		// TODO: Handle listener type does not exist
	}

	listener := listenerConstructor(listenerConfig)

	return &FilteringSpanListener{listener, filterer}
}

func crateFiltersFromConfig(filterConfigs []interface{}) []SpanFilter {
	filters := []SpanFilter{}
	for _, filterConfig := range filterConfigs {
		if filterConfig, ok := filterConfig.(map[string]interface{}); ok {
			if composite, ok := filterConfig["composite"].(bool); ok && composite {
				all, ok := filterConfig["all"].(bool)
				if !ok {
					// TODO: Handle all value is not bool
				}

				cf := &CompositeSpanFilter{
					spanFilters: []SpanFilter{},
					all:         all,
					composite:   true,
				}

				if compositeFilterConfigs, ok := filterConfig["filters"].([]interface{}); ok {
					cf.spanFilters = crateFiltersFromConfig(compositeFilterConfigs)
				}

				filters = append(filters, cf)
			} else {
				filters = append(filters, NewThundraSpanFilter(filterConfig))
			}
		}
	}

	return filters
}
