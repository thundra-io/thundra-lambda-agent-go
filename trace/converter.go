package trace

import (
	"github.com/thundra-io/thundra-lambda-agent-go/thundra_tracer"
)

func convertSpantoTraceData(trace *trace) []map[string]interface{} {
	mr := trace.GetRecorder()
	sTree := mr.GetSpanTree()
	root := walkAndConvert(sTree)
	if trace.thrownError != nil {
		root = setErrorsOnRightestChildren(root, trace.errors, trace.thrownError)
	}
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


// If panic occurs on the execution, it occurs on the last leaf node because execution can not continue from that point.
// setErrorsOnRightestChildren traverses the tree and adds error logs to all rightest children
func setErrorsOnRightestChildren(root map[string]interface{}, errors []string, thrownError interface{}) map[string]interface{} {
	aiChildren := root[auditInfoChildren].([]map[string]interface{})
	if len(aiChildren) != 0 {
		last := aiChildren[len(aiChildren)-1]
		last = setErrorsOnRightestChildren(last, errors, thrownError)
		aiChildren[len(aiChildren)-1] = last
		root[auditInfoChildren] = aiChildren
	}
	root[auditInfoErrors] = errors
	root[auditInfoThrownError] = thrownError
	return root
}
