package thundrardb

import (
	"net"
	"strings"

	"github.com/go-sql-driver/mysql"

	"github.com/thundra-io/thundra-lambda-agent-go/application"
	"github.com/thundra-io/thundra-lambda-agent-go/config"
	"github.com/thundra-io/thundra-lambda-agent-go/constants"
	"github.com/thundra-io/thundra-lambda-agent-go/tracer"
)

type mysqlIntegration struct{}

func (i *mysqlIntegration) getOperationName(query string) string {
	querySplit := strings.Split(query, " ")
	operation := ""
	if len(querySplit) > 0 {
		operation = querySplit[0]
	}
	return operation
}

func (i *mysqlIntegration) beforeCall(query string, span *tracer.RawSpan, dsn string) {
	span.ClassName = constants.ClassNames["MYSQL"]
	span.DomainName = constants.DomainNames["DB"]

	operation := i.getOperationName(query)

	dbName := ""
	host := ""
	port := ""
	cfg, err := mysql.ParseDSN(dsn)

	if err == nil {
		dbName = cfg.DBName
		host, port, err = net.SplitHostPort(cfg.Addr)
	}

	// Set span tags
	tags := map[string]interface{}{
		constants.SpanTags["OPERATION_TYPE"]:          operationToType[strings.ToLower(operation)],
		constants.SpanTags["TRIGGER_OPERATION_NAMES"]: []string{application.FunctionName},
		constants.SpanTags["TRIGGER_DOMAIN_NAME"]:     constants.AwsLambdaApplicationDomain,
		constants.SpanTags["TRIGGER_CLASS_NAME"]:      constants.AwsLambdaApplicationClass,
		constants.SpanTags["TOPOLOGY_VERTEX"]:         true,
		constants.DBTags["DB_STATEMENT_TYPE"]:         strings.ToUpper(operation),
		constants.DBTags["DB_TYPE"]:                   "mysql",
		constants.DBTags["DB_STATEMENT_TYPE"]:         strings.ToUpper(operation),
		constants.DBTags["DB_INSTANCE"]:               dbName,
		constants.DBTags["DB_HOST"]:                   host,
		constants.DBTags["DB_PORT"]:                   port,
	}

	if !config.MaskRDBStatement {
		tags[constants.DBTags["DB_STATEMENT"]] = query
	}

	span.Tags = tags
}

func (i *mysqlIntegration) afterCall(query string, span *tracer.RawSpan, dsn string) {
	return
}
