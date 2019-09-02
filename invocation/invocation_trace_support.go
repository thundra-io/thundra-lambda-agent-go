package invocation

import (
	"errors"
	"strings"

	"github.com/thundra-io/thundra-lambda-agent-go/config"
	"github.com/thundra-io/thundra-lambda-agent-go/constants"
	"github.com/thundra-io/thundra-lambda-agent-go/trace"
	"github.com/thundra-io/thundra-lambda-agent-go/tracer"
	"github.com/thundra-io/thundra-lambda-agent-go/utils"
)

var incomingTraceLinks = make([]string, 0)

// Resource type stores information about resources in spans
type Resource struct {
	ResourceType        string   `json:"resourceType"`
	ResourceName        string   `json:"resourceName"`
	ResourceOperation   string   `json:"resourceOperation"`
	ResourceCount       int      `json:"resourceCount"`
	ResourceErrorCount  int      `json:"resourceErrorCount"`
	ResourceDuration    int64    `json:"resourceDuration"`
	ResourceMaxDuration int64    `json:"resourceMaxDuration"`
	ResourceAvgDuration float64  `json:"resourceAvgDuration"`
	ResourceErrors      []string `json:"resourceErrors"`
	resourceErrorsMap   map[string]struct{}
}

func (r *Resource) accept(rawSpan *tracer.RawSpan) bool {
	if rawSpan == nil {
		return false
	}

	operationType, ok := rawSpan.GetTag(constants.SpanTags["OPERATION_TYPE"]).(string)
	if !ok {
		operationType = ""
	}

	return strings.ToUpper(r.ResourceType) == strings.ToUpper(rawSpan.ClassName) && r.ResourceName == rawSpan.OperationName && r.ResourceOperation == operationType
}

func (r *Resource) merge(rawSpan *tracer.RawSpan) {
	if !r.accept(rawSpan) || rawSpan == nil {
		return
	}
	r.ResourceCount++
	r.ResourceDuration += rawSpan.Duration()
	r.ResourceAvgDuration = utils.Round(float64(r.ResourceDuration)/float64(r.ResourceCount)*100) / 100
	if rawSpan.Duration() > r.ResourceMaxDuration {
		r.ResourceMaxDuration = rawSpan.Duration()
	}
	erroneous, ok := rawSpan.GetTag(constants.AwsError).(bool)
	if ok && erroneous {
		r.ResourceErrorCount++
		errKind, ok := rawSpan.GetTag(constants.AwsErrorKind).(string)
		if ok {
			r.resourceErrorsMap[errKind] = void
		}
	}
}

// NewResource initializes and returns Resource with given span
func NewResource(rawSpan *tracer.RawSpan) (*Resource, error) {
	if rawSpan == nil {
		return &Resource{}, errors.New("Nil span")
	}

	operationType, ok := rawSpan.GetTag(constants.SpanTags["OPERATION_TYPE"]).(string)
	if !ok {
		operationType = ""
	}
	erroneous, ok := rawSpan.GetTag(constants.AwsError).(bool)
	errorCount := 0

	resourceErrorsMap := make(map[string]struct{})
	if ok && erroneous {
		errorCount = 1
		errKind, ok := rawSpan.GetTag(constants.AwsErrorKind).(string)
		if ok {
			resourceErrorsMap[errKind] = void
		}
	}

	resource := Resource{
		ResourceType:        rawSpan.ClassName,
		ResourceName:        rawSpan.OperationName,
		ResourceOperation:   operationType,
		ResourceCount:       1,
		ResourceErrorCount:  errorCount,
		ResourceDuration:    rawSpan.Duration(),
		ResourceMaxDuration: rawSpan.Duration(),
		ResourceAvgDuration: float64(rawSpan.Duration()),
		resourceErrorsMap:   resourceErrorsMap,
	}
	return &resource, nil
}

func getResourceID(rawSpan *tracer.RawSpan) string {
	resourceID := ""
	if rawSpan == nil {
		return resourceID
	}
	operationType, ok := rawSpan.GetTag(constants.SpanTags["OPERATION_TYPE"]).(string)
	if !ok {
		operationType = ""
	}

	return strings.ToUpper(rawSpan.ClassName) + rawSpan.OperationName + operationType
}

func getResources(rootSpanID string) []Resource {
	resources := make(map[string]*Resource)
	spanList := trace.GetInstance().Recorder.GetSpans()
	for _, s := range spanList {
		vertex, ok := s.GetTag(constants.SpanTags["TOPOLOGY_VERTEX"]).(bool)
		if !ok || !vertex || s.Context.SpanID == rootSpanID {
			continue
		}
		resourceID := getResourceID(s)
		if resource, exist := resources[resourceID]; exist {
			resource.merge(s)
		} else {
			resource, err := NewResource(s)
			if err == nil {
				resources[resourceID] = resource
			}
		}
	}

	values := make([]Resource, 0, len(resources))
	for _, value := range resources {
		var resourceErrors = make([]string, 0)
		for k := range value.resourceErrorsMap {
			resourceErrors = append(resourceErrors, k)
		}
		value.ResourceErrors = resourceErrors
		values = append(values, *value)
	}
	return values
}

func getIncomingTraceLinks() []string {
	if config.ThundraDisabled {
		return []string{}
	}

	incomingTraceLinksMap := make(map[string]struct{})

	for _, link := range incomingTraceLinks {
		incomingTraceLinksMap[link] = void
	}

	links := []string{}
	for k := range incomingTraceLinksMap {
		links = append(links, k)
	}
	return links
}

func getOutgoingTraceLinks() []string {
	var outGoingTraceLinks = make([]string, 0)
	if config.ThundraDisabled {
		return outGoingTraceLinks
	}

	outgoingTraceLinksMap := make(map[string]struct{})

	spanList := trace.GetInstance().Recorder.GetSpans()

	for _, s := range spanList {
		links, ok := s.GetTag(constants.SpanTags["TRACE_LINKS"]).([]string)
		if ok {
			for _, link := range links {
				outgoingTraceLinksMap[link] = void
			}
		}
	}

	for k := range outgoingTraceLinksMap {
		outGoingTraceLinks = append(outGoingTraceLinks, k)
	}
	return outGoingTraceLinks
}

// AddIncomingTraceLinks adds links to incomingTraceLinks
func AddIncomingTraceLinks(links []string) {
	for _, link := range links {
		incomingTraceLinks = append(incomingTraceLinks, link)
	}
}

func clearTraceLinks() {
	incomingTraceLinks = []string{}
}
