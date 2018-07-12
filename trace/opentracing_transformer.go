package trace

import (
	"github.com/thundra-io/thundra-lambda-agent-go/thundra_tracer"
)

// transformSpantoTraceData transforms manually instrumented spans to traceData format
func transformSpantoTraceData(recorder thundra_tracer.SpanRecorder) interface{} {
	sTree := recorder.GetSpanTree()
	if sTree == nil {
		return nil
	}
	root := traverseAndTransform(sTree)
	recorder.Reset()
	// We currently ignore the root span because of our trace plugin sets the the root
	return root //root[auditInfoChildren]
}

// traverseAndTransform traverses tree in depth-first and transforms spans to audit data
func traverseAndTransform(t *thundra_tracer.RawSpanTree) map[string]interface{} {
	if t == nil {
		return nil
	}
	parent := transformToAuditData(t.Value)
	for _, child := range t.Children {
		ad := traverseAndTransform(child)
		// We need to convert it because it is an interface
		aiChildren, ok := parent[auditInfoChildren].([]map[string]interface{})
		if ok {
			aiChildren = append(aiChildren, ad)
			parent[auditInfoChildren] = aiChildren
		}
	}
	return parent
}

func transformToAuditData(span *thundra_tracer.RawSpan) map[string]interface{} {
	return map[string]interface{}{
		auditInfoContextName:    span.Operation,
		auditInfoId:             span.Context.SpanID,
		auditInfoOpenTimestamp:  span.StartTimestamp,
		auditInfoCloseTimestamp: span.EndTimestamp,
		auditInfoErrors:         nil,
		auditInfoThrownError:    nil,
		auditInfoChildren:       make([]map[string]interface{}, 0),
		auditInfoProps:          span.Tags,
	}
}
