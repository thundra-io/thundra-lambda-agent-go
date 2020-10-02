package thundraaws

import (
	"github.com/thundra-io/thundra-lambda-agent-go/v2/utils"

	"github.com/thundra-io/thundra-lambda-agent-go/v2/config"

	"github.com/aws/aws-sdk-go/aws/request"
	"github.com/thundra-io/thundra-lambda-agent-go/v2/application"
	"github.com/thundra-io/thundra-lambda-agent-go/v2/constants"
	"github.com/thundra-io/thundra-lambda-agent-go/v2/tracer"
)

type athenaIntegration struct {
	fields map[string]interface{}
}

func (i *athenaIntegration) parseFields(r *request.Request) {
	i.fields = utils.SerializeToMap(r.Params)
}

func (i *athenaIntegration) getOperationName(r *request.Request) string {
	i.parseFields(r)

	dbName := i.getDatabaseName()
	if len(dbName) > 0 {
		return dbName
	}

	return constants.AWSServiceRequest
}

func (i *athenaIntegration) getNestedField(firstFieldName, secondFieldName string) string {
	if firstFieldInterface, ok := i.fields[firstFieldName]; ok {
		if firstField, ok := firstFieldInterface.(map[string]interface{}); ok {
			if secondFieldInterface, ok := firstField[secondFieldName]; ok {
				if secondField, ok := secondFieldInterface.(string); ok {
					return secondField
				}
			}
		}
	}
	return ""
}

func (i *athenaIntegration) getField(fieldName string) string {
	if valInterface, ok := i.fields[fieldName]; ok {
		if val, ok := valInterface.(string); ok {
			return val
		}
	}
	return ""
}

func (i *athenaIntegration) getOutputLocation() string {
	return i.getNestedField("ResultConfiguration", "OutputLocation")
}

func (i *athenaIntegration) getQuery() string {
	return i.getField("QueryString")
}

func (i *athenaIntegration) getQueryExecutionIDs(named bool) []string {
	var singleFieldName string
	var multiFieldName string

	if named {
		singleFieldName = "NamedQueryId"
		multiFieldName = "NamedQueryIds"
	} else {
		singleFieldName = "QueryExecutionId"
		multiFieldName = "QueryExecutionIds"
	}

	if qeidInterface, ok := i.fields[singleFieldName]; ok {
		if qeid, ok := qeidInterface.(string); ok {
			return []string{qeid}
		}
	}
	if qeidsInterface, ok := i.fields[multiFieldName]; ok {
		if qeids, ok := qeidsInterface.([]interface{}); ok {
			res := []string{}
			for _, p := range qeids {
				if v, ok := p.(string); ok {
					res = append(res, v)
				}
			}
			return res
		}
	}
	return []string{}
}

func (i *athenaIntegration) getDatabaseName() string {
	dbName := i.getField("Database")
	if len(dbName) > 0 {
		return dbName
	}
	return i.getNestedField("QueryExecutionContext", "Database")
}

func (i *athenaIntegration) beforeCall(r *request.Request, span *tracer.RawSpan) {
	span.ClassName = constants.ClassNames["ATHENA"]
	span.DomainName = constants.DomainNames["DB"]

	operationName := r.Operation.Name
	operationType := getOperationType(operationName, constants.ClassNames["ATHENA"])

	tags := map[string]interface{}{
		constants.SpanTags["OPERATION_TYPE"]:          operationType,
		constants.AwsSDKTags["REQUEST_NAME"]:          operationName,
		constants.SpanTags["TOPOLOGY_VERTEX"]:         true,
		constants.SpanTags["TRIGGER_OPERATION_NAMES"]: []string{application.FunctionName},
		constants.SpanTags["TRIGGER_DOMAIN_NAME"]:     constants.AwsLambdaApplicationDomain,
		constants.SpanTags["TRIGGER_CLASS_NAME"]:      constants.AwsLambdaApplicationClass,
	}

	dbName := i.getDatabaseName()
	outputLocation := i.getOutputLocation()
	queryExecutionIDs := i.getQueryExecutionIDs(false)
	namedQueryIDs := i.getQueryExecutionIDs(true)

	if len(dbName) > 0 {
		tags[constants.DBTags["DB_INSTANCE"]] = dbName
	}

	if len(outputLocation) > 0 {
		tags[constants.AwsAthenaTags["S3_OUTPUT_LOCATION"]] = outputLocation
	}

	if len(queryExecutionIDs) > 0 {
		tags[constants.AwsAthenaTags["REQUEST_QUERY_EXECUTION_IDS"]] = queryExecutionIDs
	}

	if len(namedQueryIDs) > 0 {
		tags[constants.AwsAthenaTags["REQUEST_NAMED_QUERY_IDS"]] = namedQueryIDs
	}

	if !config.MaskAthenaStatement {
		if q := i.getQuery(); len(q) > 0 {
			tags[constants.DBTags["DB_STATEMENT"]] = q
		}
	}

	span.Tags = tags
}

func (i *athenaIntegration) afterCall(r *request.Request, span *tracer.RawSpan) {
	response := utils.SerializeToMap(r.Data)

	if qeidInterface, ok := response["QueryExecutionId"]; ok {
		if qeid, ok := qeidInterface.(string); ok {
			span.Tags[constants.AwsAthenaTags["RESPONSE_QUERY_EXECUTION_IDS"]] = []string{qeid}
		}
	}

	if qeidsInterface, ok := response["QueryExecutionIds"]; ok {
		if qeids, ok := qeidsInterface.([]string); ok {
			span.Tags[constants.AwsAthenaTags["RESPONSE_QUERY_EXECUTION_IDS"]] = qeids
		}
	}

	if nqidInterface, ok := response["NamedQueryId"]; ok {
		if nqid, ok := nqidInterface.(string); ok {
			span.Tags[constants.AwsAthenaTags["RESPONSE_NAMED_QUERY_IDS"]] = []string{nqid}
		}
	}

	if nqidsInterface, ok := response["NamedQueryIds"]; ok {
		if nqids, ok := nqidsInterface.([]string); ok {
			span.Tags[constants.AwsAthenaTags["RESPONSE_QUERY_EXECUTION_IDS"]] = nqids
		}
	}
}

func init() {
	integrations["Athena"] = &athenaIntegration{}
}
