package tracer

import (
	"fmt"

	ot "github.com/opentracing/opentracing-go"
)

type SpanFilter interface {
	Accept(*spanImpl) bool
}

type SpanFilterer interface {
	Accept(*spanImpl) bool
}

type ThundraSpanFilterer struct {
	spanFilters []SpanFilter
	all         bool
}

type ThundraSpanFilter struct {
	DomainName    string
	ClassName     string
	OperationName string
	Reverse       bool
	Tags          ot.Tags
}

type CompositeSpanFilter struct {
	spanFilters []SpanFilter
	all         bool
	composite   bool
}

func (f *CompositeSpanFilter) Accept(span *spanImpl) bool {
	res := f.all
	for _, sf := range f.spanFilters {
		if f.all {
			res = res && sf.Accept(span)
		} else {
			res = res || sf.Accept(span)
		}
	}
	return res
}

func (t *ThundraSpanFilterer) Accept(span *spanImpl) bool {
	res := t.all
	for _, sf := range t.spanFilters {
		if t.all {
			res = res && sf.Accept(span)
		} else {
			res = res || sf.Accept(span)
		}
	}
	return res
}

func (t *ThundraSpanFilterer) AddFilter(sf SpanFilter) {
	t.spanFilters = append(t.spanFilters, sf)
}

func (t *ThundraSpanFilterer) ClearFilters() {
	t.spanFilters = []SpanFilter{}
}

func (t *ThundraSpanFilter) Accept(span *spanImpl) bool {
	accepted := true
	if span == nil {
		return accepted
	}

	if t.DomainName != "" {
		accepted = (t.DomainName == span.raw.DomainName)
	}

	if accepted && t.ClassName != "" {
		accepted = (t.ClassName == span.raw.ClassName)
	}

	if accepted && t.OperationName != "" {
		accepted = (t.OperationName == span.raw.OperationName)
	}

	if accepted && t.Tags != nil {
		for k, v := range t.Tags {
			if fmt.Sprintf("%v", span.raw.GetTag(k)) != fmt.Sprintf("%v", v) {
				accepted = false
				break
			}
		}
	}

	if t.Reverse {
		return !accepted
	}

	return accepted
}

// NewThundraSpanFilter creates and returns a new ThundraSpanFilter from config
func NewThundraSpanFilter(config map[string]interface{}) *ThundraSpanFilter {
	spanFilter := ThundraSpanFilter{}

	if domainName, ok := config["domainName"].(string); ok {
		spanFilter.DomainName = domainName
	}
	if className, ok := config["className"].(string); ok {
		spanFilter.ClassName = className
	}
	if operationName, ok := config["operationName"].(string); ok {
		spanFilter.OperationName = operationName
	}
	if reverse, ok := config["reverse"].(bool); ok {
		spanFilter.Reverse = reverse
	}
	if tags, ok := config["tags"].(map[string]interface{}); ok {
		spanFilter.Tags = tags
	}

	return &spanFilter
}
