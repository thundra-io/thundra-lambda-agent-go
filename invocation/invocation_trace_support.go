package invocation

import (
	"github.com/thundra-io/thundra-lambda-agent-go/config"
	"github.com/thundra-io/thundra-lambda-agent-go/constants"
	"github.com/thundra-io/thundra-lambda-agent-go/trace"
)

var incomingTraceLinks = make([]string, 0)

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
