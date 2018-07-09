package trace

import (
	"github.com/thundra-io/thundra-lambda-agent-go/thundra_tracer"
)

func convertSpantoTraceData(trace *trace) []map[string]interface{} {
	mr := trace.GetRecorder()
	sTree := mr.GetSpanTree()
	root := walkAndConvert(sTree)
	mr.Reset()
	aiChildren := root[auditInfoChildren].([]map[string]interface{})
	return aiChildren
}

// walkAndConvert traverses a tree depth-first and converts span data to audit data
func walkAndConvert(t *thundra_tracer.RawSpanTree) map[string]interface{} {
	if t == nil {
		return nil
	}
	parent := convertToAuditData(t.Value)
	for _, child := range t.Children {
		ad := walkAndConvert(child)
		// We need to convert it because it is an interface
		aiChildren, ok := parent[auditInfoChildren].([]map[string]interface{})
		if ok {
			aiChildren = append(aiChildren, ad)
			parent[auditInfoChildren] = aiChildren
		}
	}
	return parent
}

func convertToAuditData(span *thundra_tracer.RawSpan) map[string]interface{} {
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
