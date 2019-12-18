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
	filterer := &ThundraSpanFilterer{spanFilters: []SpanFilter{}}

	listenerDef, ok := config["listener"].(map[string]interface{})
	if !ok {
		log.Println("Listener configuration is not valid for FilteringSpanListener")
		return nil
	}

	listenerName, ok := listenerDef["type"].(string)
	listenerConstructor, ok := SpanListenerConstructorMap[listenerName]
	if !ok {
		log.Println("Given listener type is not valid for FilteringSpanListener")
		return nil
	}

	listenerConfig, ok := listenerDef["config"].(map[string]interface{})
	listener := listenerConstructor(listenerConfig)

	if all, ok := config["all"].(bool); ok {
		filterer.all = all
	}

	if filterConfigs, ok := config["filters"].([]interface{}); ok {
		filterer.spanFilters = crateFiltersFromConfig(filterConfigs)
	}

	return &FilteringSpanListener{listener, filterer}
}

func crateFiltersFromConfig(filterConfigs []interface{}) []SpanFilter {
	filters := []SpanFilter{}
	for _, filterConfig := range filterConfigs {
		if filterConfig, ok := filterConfig.(map[string]interface{}); ok {
			if composite, ok := filterConfig["composite"].(bool); ok && composite {
				cf := &CompositeSpanFilter{
					spanFilters: []SpanFilter{},
					all:         false,
					composite:   true,
				}

				if all, ok := filterConfig["all"].(bool); ok {
					cf.all = all
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
