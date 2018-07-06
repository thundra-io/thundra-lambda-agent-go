package trace

import (
	"github.com/thundra-io/thundra-lambda-agent-go/otTracer"
)

func convertSpantoTraceData(trace *trace) []map[string]interface{} {
	mr := trace.GetMemRecorder()
	sTree := mr.GetAllSpansTree()
	root := walkAndConvert(sTree)
	mr.Reset()
	aiChildren := root[auditInfoChildren].([]map[string]interface{})
	return aiChildren
}

// walkAndConvert traverses a tree depth-first
func walkAndConvert(t *otTracer.RawSpanTree) map[string]interface{} {
	if t == nil {
		return nil
	}
	root := convertToAuditData(t.Value)
	for _, child := range t.Children {
		ad := walkAndConvert(child)
		aiChildren := root[auditInfoChildren].([]map[string]interface{})
		aiChildren = append(aiChildren, ad)
		root[auditInfoChildren] = aiChildren
	}
	return root
}

func convertToAuditData(span *otTracer.RawSpan) map[string]interface{} {
	audit := map[string]interface{}{
		auditInfoContextName:    span.Operation,
		auditInfoId:             span.Context.SpanID,
		auditInfoOpenTimestamp:  span.StartTimestamp,
		auditInfoCloseTimestamp: span.EndTimestamp,
		auditInfoErrors:         nil,
		auditInfoThrownError:    nil,
		auditInfoChildren:       make([]map[string]interface{}, 0),
		auditInfoProps:          span.Tags,
	}
	return audit
}
