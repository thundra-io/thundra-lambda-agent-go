package trace


// transformSpantoTraceData transforms manually instrumented spans to traceData format
// func transformSpantoTraceData(recorder ttracer.SpanRecorder) []map[string]interface{} {
// 	sTree := recorder.GetSpanTree()
// 	if sTree == nil {
// 		return nil
// 	}
// 	root := traverseAndTransform(sTree)
// 	recorder.Reset()
// 	// Creates a new array and add the object because it should return an array
// 	var rslt []map[string]interface{}
// 	rslt = append(rslt, root)
// 	return rslt
// }

// // traverseAndTransform traverses tree in depth-first and transforms spans to audit data
// func traverseAndTransform(t *ttracer.RawSpanTree) map[string]interface{} {
// 	if t == nil {
// 		return nil
// 	}
// 	parent := transformToAuditData(t.Value)
// 	for _, child := range t.Children {
// 		ad := traverseAndTransform(child)
// 		// We need to convert it because it is an interface
// 		aiChildren, ok := parent[auditInfoChildren].([]map[string]interface{})
// 		if ok {
// 			aiChildren = append(aiChildren, ad)
// 			parent[auditInfoChildren] = aiChildren
// 		}
// 	}
// 	return parent
// }

// func transformToAuditData(span *ttracer.RawSpan) map[string]interface{} {
// 	return map[string]interface{}{
// 		auditInfoContextName:    span.Operation,
// 		auditInfoId:             span.Context.SpanID,
// 		auditInfoOpenTimestamp:  span.StartTimestamp,
// 		auditInfoCloseTimestamp: span.EndTimestamp,
// 		auditInfoErrors:         nil,
// 		auditInfoThrownError:    nil,
// 		auditInfoChildren:       make([]map[string]interface{}, 0),
// 		auditInfoProps:          span.Tags,
// 	}
// }
