package thundraaws

import (
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"reflect"
	"strconv"
	"time"

	"github.com/aws/aws-sdk-go/aws/request"
	"github.com/aws/aws-sdk-go/service/firehose"
	"github.com/thundra-io/thundra-lambda-agent-go/application"
	"github.com/thundra-io/thundra-lambda-agent-go/constants"
	"github.com/thundra-io/thundra-lambda-agent-go/tracer"
)

type firehoseIntegration struct{}

func (i *firehoseIntegration) getDeliveryStreamName(r *request.Request) string {
	fields := struct {
		DeliveryStreamName string
	}{}
	m, err := json.Marshal(r.Params)
	if err != nil {
		return ""
	}
	if err = json.Unmarshal(m, &fields); err != nil {
		return ""
	}
	if len(fields.DeliveryStreamName) > 0 {
		return fields.DeliveryStreamName
	}
	return ""
}

func (i *firehoseIntegration) getOperationName(r *request.Request) string {
	dsn := i.getDeliveryStreamName(r)
	if len(dsn) > 0 {
		return dsn
	}
	return constants.AWSServiceRequest
}

func (i *firehoseIntegration) beforeCall(r *request.Request, span *tracer.RawSpan) {
	span.ClassName = constants.ClassNames["FIREHOSE"]
	span.DomainName = constants.DomainNames["STREAM"]

	operationName := r.Operation.Name
	operationType := constants.FirehoseRequestTypes[operationName]

	tags := map[string]interface{}{
		constants.AwsFirehoseTags["STREAM_NAME"]:      i.getDeliveryStreamName(r),
		constants.SpanTags["OPERATION_TYPE"]:          operationType,
		constants.AwsSDKTags["REQUEST_NAME"]:          operationName,
		constants.SpanTags["TOPOLOGY_VERTEX"]:         true,
		constants.SpanTags["TRIGGER_OPERATION_NAMES"]: []string{application.FunctionName},
		constants.SpanTags["TRIGGER_DOMAIN_NAME"]:     constants.AwsLambdaApplicationDomain,
		constants.SpanTags["TRIGGER_CLASS_NAME"]:      constants.AwsLambdaApplicationClass,
	}

	span.Tags = tags
}

func (i *firehoseIntegration) afterCall(r *request.Request, span *tracer.RawSpan) {

	traceLinks := i.getTraceLinks(r.Operation.Name, r)
	if traceLinks != nil {
		span.Tags[constants.SpanTags["TRACE_LINKS"]] = traceLinks
	}
}

func getTimeStamp(dateStr string) int64 {
	date, err := time.Parse(time.RFC1123, dateStr)
	timestamp := time.Now().Unix() - 1
	if err == nil {
		timestamp = date.Unix()
	}
	return timestamp
}

func (i *firehoseIntegration) getTraceLinks(operationName string, r *request.Request) []string {
	requestValue := reflect.ValueOf(r.Params)
	if requestValue == (reflect.Value{}) {
		return nil
	}

	streamName := i.getDeliveryStreamName(r)
	region := ""
	dateStr := ""

	if r.Config.Region != nil {
		region = *r.Config.Region
	}
	if r.HTTPResponse != nil && r.HTTPResponse.Header != nil {
		dateStr = r.HTTPResponse.Header.Get("date")
	}

	if operationName == "PutRecord" {
		if recordInput, ok := requestValue.Elem().Interface().(firehose.PutRecordInput); ok {
			if recordInput.Record != nil {
				data := recordInput.Record.Data
				return i.generateTraceLinks(region, dateStr, data, streamName)
			}
		}
	} else if operationName == "PutRecordBatch" {
		records := requestValue.Elem().FieldByName("Records")
		if records != (reflect.Value{}) {
			var links []string
			for j := 0; j < records.Len(); j++ {
				if record, ok := records.Index(j).Elem().Interface().(firehose.Record); ok {
					linksForRecord := i.generateTraceLinks(region, dateStr, record.Data, streamName)
					for _, link := range linksForRecord {
						links = append(links, link)
					}
				}
			}
			return links
		}
	}
	return nil
}

func (i *firehoseIntegration) generateTraceLinks(region string, dateStr string, data []byte, streamName string) []string {
	var traceLinks []string
	timestamp := getTimeStamp(dateStr)

	b := md5.Sum(data)
	dataMD5 := hex.EncodeToString(b[:])

	for j := 0; j < 3; j++ {
		traceLinks = append(traceLinks, region+":"+streamName+":"+strconv.FormatInt(timestamp+int64(j), 10)+":"+dataMD5)
	}

	return traceLinks
}

func init() {
	integrations["Firehose"] = &firehoseIntegration{}
}
