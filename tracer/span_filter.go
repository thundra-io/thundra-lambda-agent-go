package tracer

import (
	"fmt"
	"strings"

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
}

type ThundraSpanFilter struct {
	DomainName    string
	ClassName     string
	OperationName string
	Tags          ot.Tags
}

func (t *ThundraSpanFilterer) Accept(span *spanImpl) bool {
	for _, sf := range t.spanFilters {
		if sf.Accept(span) {
			return true
		}
	}
	return false
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
	return accepted
}

// NewThundraSpanFilter creates and returns a new ThundraSpanFilter from config
func NewThundraSpanFilter(config map[string]string) *ThundraSpanFilter {

	spanFilter := ThundraSpanFilter{}

	if config["domainName"] != "" {
		spanFilter.DomainName = config["domainName"]
	}
	if config["className"] != "" {
		spanFilter.ClassName = config["className"]
	}
	if config["operationName"] != "" {
		spanFilter.OperationName = config["operationName"]
	}

	tagPrefix := "tag."
	tags := make(map[string]interface{})
	for k, v := range config {
		if strings.HasPrefix(k, tagPrefix) {
			tags[k[len(tagPrefix):]] = v
		}
	}
	spanFilter.Tags = tags

	return &spanFilter
}
