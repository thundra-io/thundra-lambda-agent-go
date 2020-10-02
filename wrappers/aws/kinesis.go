package thundraaws

import (
	"encoding/json"
	"reflect"

	"github.com/aws/aws-sdk-go/aws/request"
	"github.com/thundra-io/thundra-lambda-agent-go/v2/application"
	"github.com/thundra-io/thundra-lambda-agent-go/v2/constants"
	"github.com/thundra-io/thundra-lambda-agent-go/v2/tracer"
	"github.com/thundra-io/thundra-lambda-agent-go/v2/utils"
)

type kinesisIntegration struct{}

func (i *kinesisIntegration) getStreamName(r *request.Request) string {
	fields := struct {
		StreamName string
	}{}
	m, err := json.Marshal(r.Params)
	if err != nil {
		return ""
	}
	if err = json.Unmarshal(m, &fields); err != nil {
		return ""
	}
	if len(fields.StreamName) > 0 {
		return fields.StreamName
	}
	return ""
}

func (i *kinesisIntegration) getOperationName(r *request.Request) string {
	streamName := i.getStreamName(r)
	if len(streamName) > 0 {
		return streamName
	}
	return constants.AWSServiceRequest
}

func (i *kinesisIntegration) beforeCall(r *request.Request, span *tracer.RawSpan) {
	span.ClassName = constants.ClassNames["KINESIS"]
	span.DomainName = constants.DomainNames["STREAM"]

	operationName := r.Operation.Name
	operationType := getOperationType(operationName, constants.ClassNames["KINESIS"])

	tags := map[string]interface{}{
		constants.AwsKinesisTags["STREAM_NAME"]:       i.getStreamName(r),
		constants.SpanTags["OPERATION_TYPE"]:          operationType,
		constants.AwsSDKTags["REQUEST_NAME"]:          operationName,
		constants.SpanTags["TOPOLOGY_VERTEX"]:         true,
		constants.SpanTags["TRIGGER_OPERATION_NAMES"]: []string{application.FunctionName},
		constants.SpanTags["TRIGGER_DOMAIN_NAME"]:     constants.AwsLambdaApplicationDomain,
		constants.SpanTags["TRIGGER_CLASS_NAME"]:      constants.AwsLambdaApplicationClass,
	}

	span.Tags = tags
}

func (i *kinesisIntegration) afterCall(r *request.Request, span *tracer.RawSpan) {
	traceLinks := i.getTraceLinks(r)
	if traceLinks != nil {
		span.Tags[constants.SpanTags["TRACE_LINKS"]] = traceLinks
	}
}

func (i *kinesisIntegration) getTraceLinks(r *request.Request) []string {
	responseValue := reflect.ValueOf(r.Data)

	if responseValue == (reflect.Value{}) {
		return nil
	}

	records := responseValue.Elem().FieldByName("Records")
	region := ""
	streamName := i.getStreamName(r)

	if r.Config.Region != nil {
		region = *r.Config.Region
	}

	if records != (reflect.Value{}) {
		var links []string
		for j := 0; j < records.Len(); j++ {
			record := records.Index(j).Elem()
			if sequenceNumber, ok := utils.GetStringFieldFromValue(record, "SequenceNumber"); ok {
				if shardID, ok := utils.GetStringFieldFromValue(record, "ShardId"); ok {
					links = append(links, region+":"+streamName+":"+shardID+":"+sequenceNumber)
				}
			}
		}
		return links
	}
	if sequenceNumber, ok := utils.GetStringFieldFromValue(responseValue.Elem(), "SequenceNumber"); ok {
		if shardID, ok := utils.GetStringFieldFromValue(responseValue.Elem(), "ShardId"); ok {
			return []string{region + ":" + streamName + ":" + shardID + ":" + sequenceNumber}
		}
	}

	return nil
}

func init() {
	integrations["Kinesis"] = &kinesisIntegration{}
}
