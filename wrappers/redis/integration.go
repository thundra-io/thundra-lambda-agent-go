package thundraredis

import (
	"github.com/thundra-io/thundra-lambda-agent-go/application"
	"github.com/thundra-io/thundra-lambda-agent-go/config"
	"github.com/thundra-io/thundra-lambda-agent-go/constants"
	"github.com/thundra-io/thundra-lambda-agent-go/tracer"
)

func beforeCall(span *tracer.RawSpan, host string, port string, commandName string, command string) {
	span.ClassName = constants.ClassNames["REDIS"]
	span.DomainName = constants.DomainNames["CACHE"]

	// Set span tags
	tags := map[string]interface{}{
		constants.SpanTags["OPERATION_TYPE"]:      constants.RedisCommandTypes[commandName],
		constants.DBTags["DB_INSTANCE"]:           host,
		constants.DBTags["DB_STATEMENT_TYPE"]:     constants.RedisCommandTypes[commandName],
		constants.DBTags["DB_TYPE"]:               "redis",
		constants.RedisTags["REDIS_HOST"]:         host,
		constants.RedisTags["REDIS_COMMAND_TYPE"]: commandName,
	}

	tags[constants.SpanTags["TRIGGER_DOMAIN_NAME"]] = constants.AwsLambdaApplicationDomain
	tags[constants.SpanTags["TRIGGER_CLASS_NAME"]] = constants.AwsLambdaApplicationClass
	tags[constants.SpanTags["TRIGGER_OPERATION_NAMES"]] = []string{application.FunctionName}
	tags[constants.SpanTags["TOPOLOGY_VERTEX"]] = true

	if !config.MaskRedisCommand {
		tags[constants.DBTags["DB_STATEMENT"]] = command
		tags[constants.RedisTags["REDIS_COMMAND"]] = command
	}

	span.Tags = tags
}
