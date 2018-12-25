package trace

import (
	"github.com/thundra-io/thundra-lambda-agent-go/plugin"
	"github.com/thundra-io/thundra-lambda-agent-go/thundra"
	"github.com/thundra-io/thundra-lambda-agent-go/thundra_tracer"
)

type spanDataStruct struct {//Base fields
	Id                        string                 `json:"id"`
	Type                      string                 `json:"type"`
	AgentVersion              string                 `json:"agentVersion"`
	DataModelVersion          string                 `json:"dataModelVersion"`
	ApplicationId             string                 `json:"applicationId"`
	ApplicationDomainName     string                 `json:"applicationDomainName"`
	ApplicationClassName      string                 `json:"applicationClassName"`
	ApplicationName           string                 `json:"applicationName"`
	ApplicationVersion        string                 `json:"applicationVersion"`
	ApplicationStage          string                 `json:"applicationStage"`
	ApplicationRuntime        string                 `json:"applicationRuntime"`
	ApplicationRuntimeVersion string                 `json:"applicationRuntimeVersion"`
	ApplicationTags           map[string]interface{} `json:"applicationTags"`

	TraceId          string                 		`json:"traceId"`
	TransactionId    string							`json:"transactionId"`
	SpanOrder		 int64							`json:"spanOrder"`
	DomainName		 string							`json:"domainName"`
	ClassName		 string							`json:"className"`
	ServiceName		 string							`json:"serviceName"`
	OperationName	 string							`json:"operationName"`
	StartTimeStamp   int64							`json:"startTimeStamp"`
	FinishTimeStamp  int64							`json:"finishTimeStamp"`
	Tags 			 map[string]interface{} 		`json:"tags"`
	Logs 			 map[string]interface{} 		`json:"logs"`

}

func prepareSpanData(span thundra_tracer.RawSpan) spanDataStruct{
	return spanDataStruct{
		Id:                        plugin.GenerateNewId(),
		Type:                      "Span",
		AgentVersion:              plugin.AgentVersion,
		DataModelVersion:          plugin.DataModelVersion,
		ApplicationId:             plugin.ApplicationId,
		ApplicationDomainName:     plugin.ApplicationDomainName,
		ApplicationClassName:      plugin.ApplicationClassName,
		ApplicationName:           plugin.FunctionName,
		ApplicationVersion:        plugin.ApplicationVersion,
		ApplicationStage:          plugin.ApplicationStage,
		ApplicationRuntime:        plugin.ApplicationRuntime,
		ApplicationRuntimeVersion: plugin.ApplicationRuntimeVersion,
		ApplicationTags:           map[string]interface{}{}, // empty object

		TraceId: "",
		TransactionId: "",
	}
}

type wrappedSpanData struct{
	ApiKey 				string 					`json:"apiKey"`
	Type 				string 					`json:"type"`
	DataModelVersion 	string 					`json:"dataModelVersion"`
	Data 				spanDataStruct			`json:"data"`
}

func wrap(data spanDataStruct) wrappedSpanData {
	return wrappedSpanData{
		ApiKey: plugin.ApiKey,
		Type: "Span",
		DataModelVersion: thundra.DataModelVersion,
		Data: data,
	}
}

// transformSpantoTraceData transforms manually instrumented spans to traceData format
func transformSpantoTraceData(recorder thundra_tracer.SpanRecorder) []wrappedSpanData {
	var spanStack []thundra_tracer.RawSpan
	//if sTree == nil {
	//	return nil
	//}
	//root := traverseAndTransform(sTree)
	//recorder.Reset()
	//// Creates a new array and add the object because it should return an array
	//var rslt []map[string]interface{}
	//rslt = append(rslt, root)
	//return rslt
	var spanDataList []wrappedSpanData
	for _, span :=  range spanStack{
		spanData := prepareSpanData(span)
		wrappedData := wrap(spanData)
		spanDataList = append(spanDataList, wrappedData)
	}
	recorder.Reset()
	return spanDataList
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
